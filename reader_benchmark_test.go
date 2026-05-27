package gocsv

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const (
	benchmarkSmallRows = 100
	benchmarkLargeRows = 100_000
	benchmarkColumns   = 8
)

type benchmarkRecord struct {
	ID        uint       `csv:"id"`
	Count     int        `csv:"count"`
	Name      string     `csv:"name"`
	Active    bool       `csv:"active"`
	Score     float64    `csv:"score"`
	CreatedAt time.Time  `csv:"created_at"`
	Note      *string    `csv:"note"`
	EndedAt   *time.Time `csv:"ended_at"`
}

// BenchmarkReadFileSmall parses a 100-row, 8-column CSV generated before the timer starts.
func BenchmarkReadFileSmall(b *testing.B) {
	benchmarkReadFileRows(b, benchmarkSmallRows)
}

// BenchmarkReadFileLarge parses a 100,000-row, 8-column CSV generated before the timer starts.
func BenchmarkReadFileLarge(b *testing.B) {
	benchmarkReadFileRows(b, benchmarkLargeRows)
}

func benchmarkReadFileRows(b *testing.B, rowCount int) {
	b.Helper()

	path, fileSize := writeBenchmarkCSV(b, rowCount)

	b.ReportAllocs()
	b.SetBytes(fileSize)
	b.Logf(
		"input=%d rows, %d columns, repeated dummy row; temp CSV generated before timer; size=%d bytes; go=%s; os=%s; arch=%s; num_cpu=%d; gomaxprocs=%d",
		rowCount,
		benchmarkColumns,
		fileSize,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.GOMAXPROCS(0),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var records []benchmarkRecord
		if err := ReadFile(&records, path, ','); err != nil {
			b.Fatal(err)
		}
		if len(records) != rowCount {
			b.Fatalf("got %d records, want %d", len(records), rowCount)
		}
		if rowCount > 0 && (records[0].ID != 1 || records[rowCount-1].EndedAt == nil) {
			b.Fatalf("records parsed incorrectly: first=%+v last=%+v", records[0], records[rowCount-1])
		}
		runtime.KeepAlive(records)
	}

	b.ReportMetric(float64(rowCount), "rows/op")
	b.ReportMetric(float64(benchmarkColumns), "cols/op")
	b.ReportMetric(float64(runtime.NumCPU()), "cpus")
	b.ReportMetric(float64(runtime.GOMAXPROCS(0)), "gomaxprocs")
}

func writeBenchmarkCSV(b *testing.B, rowCount int) (string, int64) {
	b.Helper()

	path := filepath.Join(b.TempDir(), "benchmark.csv")
	data := "id,count,name,active,score,created_at,note,ended_at\n" +
		strings.Repeat("1,-7,Alice,true,98.5,2024-01-10,first,2024-01-11\n", rowCount)

	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		b.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		b.Fatal(err)
	}

	return path, info.Size()
}
