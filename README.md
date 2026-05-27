# gocsv

GoLang package for parsing `.csv` files into structure with defined header
names in the tag.

### Example

```go
type Person struct {
	ID        uint      `csv:"ID"`
	FirstName string    `csv:"FirstName"`
	LastName  string    `csv:"LastName"`
	Age       uint      `csv:"Age"`
	Country   string    `csv:"Country"`
	Gender    string    `csv:"Gender"`
	Phone     string    `csv:"Phone"`
	Street    string    `csv:"Street"`
	Email     string    `csv:"Email"`
	Language  string    `csv:"Language"`
	BirthDate time.Time `csv:"BirthDate"`
	Added     time.Time `csv:"Added"`
	Modified  time.Time `csv:"Modified"`
}

func main() {
	persons := make([]Person, 0)

	err := gocsv.ReadFile(&persons, "person.csv", ';')
	if err != nil {
		log.Fatal(err)
	}

	for _, person := range persons {
		fmt.Println(person)
	}
}
```

### Public API

Most code only needs `ReadFile`. The `Reader` API is for cases where you
need to change CSV parser settings before reading.

#### Simple usage

```go
// ReadFile opens filePath, reads all CSV rows, and appends them to target.
// Use this for normal CSV imports.
func ReadFile(target any, filePath string, separator rune) error
```

#### Advanced usage

```go
// NewReader opens filePath and returns a configurable Reader.
// Use this when you need options such as Limit, LazyQuotes, or ApplyColumnParser.
func NewReader(filePath string, separator rune) (*Reader, error)
```

```go
// Reader holds the CSV file and parser settings.
// Create it with NewReader, configure it, then call Read.
type Reader struct { ... }
```

```go
// Read parses CSV rows into target.
// target must be a pointer to a slice of structs.
func (r *Reader) Read(target any) error
```

```go
// Close closes the opened CSV file.
// Call it when using NewReader directly, usually with defer.
func (r *Reader) Close() error
```

```go
// Limit stops reading after the configured row limit.
// Use it for previews, tests, or partial imports.
func (r *Reader) Limit(limit int) *Reader
```

```go
// ApplyColumnParser changes CSV values before they are assigned to struct fields.
// If applyToHeader is true, it also changes header names before matching csv tags.
func (r *Reader) ApplyColumnParser(parser ColumnParser, applyToHeader bool) *Reader
```

```go
// ColumnParser receives one CSV value and returns the value that should be used.
// Use it with ApplyColumnParser for normalization such as trimming or header mapping.
type ColumnParser func(value string) (string, error)
```

```go
// TrimUTF8Leading controls whether a UTF-8 BOM is removed from the first header.
// The default is true.
func (r *Reader) TrimUTF8Leading(trim bool) *Reader
```

```go
// TrimLeadingSpace controls whether leading spaces are removed from unquoted fields.
// It passes through to encoding/csv.Reader.TrimLeadingSpace.
func (r *Reader) TrimLeadingSpace(trim bool) *Reader
```

```go
// LazyQuotes controls whether malformed quoted CSV fields are accepted.
// It passes through to encoding/csv.Reader.LazyQuotes.
func (r *Reader) LazyQuotes(enable bool) *Reader
```
