package gocsv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type productWithEndDate struct {
	ID            uint       `csv:"id"`
	Name          string     `csv:"name"`
	DeactivatedAt *time.Time `csv:"deactivated_at"`
}

func TestReadFileEmptyCellDoesNotReusePreviousRowValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "products.csv")
	data := []byte("id,name,deactivated_at\n" +
		"1,Old product,2024-01-10\n" +
		"2,Active product,\n" +
		"3,Another active product,\n")

	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

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
