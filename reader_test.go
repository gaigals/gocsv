package gocsv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type basicRecord struct {
	ID        uint       `csv:"id"`
	Count     int        `csv:"count"`
	Name      string     `csv:"name"`
	Active    bool       `csv:"active"`
	Score     float64    `csv:"score"`
	CreatedAt time.Time  `csv:"created_at"`
	Note      *string    `csv:"note"`
	EndedAt   *time.Time `csv:"ended_at"`
}

type productWithEndDate struct {
	ID            uint       `csv:"id"`
	Name          string     `csv:"name"`
	DeactivatedAt *time.Time `csv:"deactivated_at"`
}

type embeddedContact struct {
	Phone string `csv:"phone"`
	Email string `csv:"email"`
}

type embeddedCustomer struct {
	Name    string          `csv:"name"`
	Contact embeddedContact `csv:"embedded"`
}

type parsedHeaderRecord struct {
	ID   uint   `csv:"id"`
	Name string `csv:"name"`
}

func writeTestCSV(t *testing.T, name string, data string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestReadFileParsesBasicFieldTypes(t *testing.T) {
	path := writeTestCSV(t, "basic.csv", "id,count,name,active,score,created_at,note,ended_at\n"+
		"1,-7,Alice,true,98.5,2024-01-10,first,2024-01-11\n"+
		"2,3,Bob,false,12.25,2024-02-20,,\n")

	var records []basicRecord
	if err := ReadFile(&records, path, ','); err != nil {
		t.Fatal(err)
	}

	if len(records) != 2 {
		t.Fatalf("got %d records, want 2", len(records))
	}
	if records[0].ID != 1 ||
		records[0].Count != -7 ||
		records[0].Name != "Alice" ||
		!records[0].Active ||
		records[0].Score != 98.5 {
		t.Fatalf("row 1 parsed incorrectly: %+v", records[0])
	}
	if records[0].CreatedAt.Format(time.DateOnly) != "2024-01-10" {
		t.Fatalf("row 1 created_at = %s, want 2024-01-10", records[0].CreatedAt.Format(time.DateOnly))
	}
	if records[0].Note == nil || *records[0].Note != "first" {
		t.Fatalf("row 1 note = %v, want first", records[0].Note)
	}
	if records[0].EndedAt == nil || records[0].EndedAt.Format(time.DateOnly) != "2024-01-11" {
		t.Fatalf("row 1 ended_at = %v, want 2024-01-11", records[0].EndedAt)
	}
	if records[1].ID != 2 ||
		records[1].Count != 3 ||
		records[1].Name != "Bob" ||
		records[1].Active ||
		records[1].Score != 12.25 {
		t.Fatalf("row 2 parsed incorrectly: %+v", records[1])
	}
	if records[1].CreatedAt.Format(time.DateOnly) != "2024-02-20" {
		t.Fatalf("row 2 created_at = %s, want 2024-02-20", records[1].CreatedAt.Format(time.DateOnly))
	}
	if records[1].Note != nil {
		t.Fatalf("row 2 note = %v, want nil", records[1].Note)
	}
	if records[1].EndedAt != nil {
		t.Fatalf("row 2 ended_at = %v, want nil", records[1].EndedAt)
	}
}

func TestReadFileParsesEmbeddedStructFields(t *testing.T) {
	path := writeTestCSV(t, "embedded.csv", "name,phone,email\n"+
		"Alice,111,alice@example.test\n"+
		"Bob,222,bob@example.test\n")

	var customers []embeddedCustomer
	if err := ReadFile(&customers, path, ','); err != nil {
		t.Fatal(err)
	}

	if len(customers) != 2 {
		t.Fatalf("got %d customers, want 2", len(customers))
	}
	if customers[0].Name != "Alice" ||
		customers[0].Contact.Phone != "111" ||
		customers[0].Contact.Email != "alice@example.test" {
		t.Fatalf("row 1 parsed incorrectly: %+v", customers[0])
	}
	if customers[1].Name != "Bob" ||
		customers[1].Contact.Phone != "222" ||
		customers[1].Contact.Email != "bob@example.test" {
		t.Fatalf("row 2 parsed incorrectly: %+v", customers[1])
	}
}

func TestReaderApplyColumnParserToHeader(t *testing.T) {
	path := writeTestCSV(t, "headers.csv", "ID,NAME\n1,Alice\n")

	reader, err := NewReader(path, ',')
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	reader.ApplyColumnParser(func(value string) (string, error) {
		switch value {
		case "ID":
			return "id", nil
		case "NAME":
			return "name", nil
		default:
			return value, nil
		}
	}, true)

	var records []parsedHeaderRecord
	if err := reader.Read(&records); err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}
	if records[0].ID != 1 || records[0].Name != "Alice" {
		t.Fatalf("row parsed incorrectly: %+v", records[0])
	}
}

func TestReadFileUsesConfiguredSeparator(t *testing.T) {
	path := writeTestCSV(t, "semicolon.csv", "id;name\n1;Alice\n")

	var records []parsedHeaderRecord
	if err := ReadFile(&records, path, ';'); err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatalf("got %d records, want 1", len(records))
	}
	if records[0].ID != 1 || records[0].Name != "Alice" {
		t.Fatalf("row parsed incorrectly: %+v", records[0])
	}
}

func TestReadFileEmptyCellDoesNotReusePreviousRowValue(t *testing.T) {
	path := writeTestCSV(t, "products.csv", "id,name,deactivated_at\n"+
		"1,Old product,2024-01-10\n"+
		"2,Active product,\n"+
		"3,Another active product,\n")

	var products []productWithEndDate
	if err := ReadFile(&products, path, ','); err != nil {
		t.Fatal(err)
	}

	if len(products) != 3 {
		t.Fatalf("got %d products, want 3", len(products))
	}
	if products[0].DeactivatedAt == nil {
		t.Fatal("row 1 deactivated_at is nil, want parsed time")
	}
	if products[1].DeactivatedAt != nil {
		t.Fatalf("row 2 deactivated_at = %v, want nil", products[1].DeactivatedAt)
	}
	if products[2].DeactivatedAt != nil {
		t.Fatalf("row 3 deactivated_at = %v, want nil", products[2].DeactivatedAt)
	}
}
