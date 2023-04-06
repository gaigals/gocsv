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
