package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gaigals/gocsv"
)

type ContactInfo struct {
	Phone  string `csv:"Phone"`
	Street string `csv:"Street"`
	Email  string `csv:"Email"`
}

type Person struct {
	ID        uint      `csv:"ID"`
	FirstName string    `csv:"FirstName"`
	LastName  string    `csv:"LastName"`
	Age       uint      `csv:"Age"`
	Country   string    `csv:"Country"`
	Gender    string    `csv:"Gender"`
	Language  string    `csv:"Language"`
	BirthDate time.Time `csv:"BirthDate"`
	Added     time.Time `csv:"Added"`
	Modified  time.Time `csv:"Modified"`

	Contacts ContactInfo `csv:"embedded"`
}

func main() {
	persons := make([]Person, 0)

	err := gocsv.ReadFile(&persons, "../test_files/person.csv", ';')
	if err != nil {
		log.Fatal(err)
	}

	for _, person := range persons {
		fmt.Println(person)
	}
}
