package gocsv

import (
	"log"
	"testing"
	"time"
)

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

func BenchmarkReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		persons := make([]Person, 0)

		reader, err := NewReader("test_files/r20c10.csv", ';')
		if err != nil {
			log.Fatalln(err)
		}

		err = reader.Read(&persons)
		if err != nil {
			log.Fatalln(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkReadFileChan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		persons := make([]Person, 0)

		reader, err := NewReader("test_files/r20c10.csv", ';')
		if err != nil {
			log.Fatalln(err)
		}

		err = reader.ReadChan(&persons)
		if err != nil {
			log.Fatalln(err)
		}
	}

	b.ReportAllocs()
}
