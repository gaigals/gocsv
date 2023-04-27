package gocsv

type columns []column

func (cols columns) maxIndex() int {
	max := 0

	for _, col := range cols {
		if col.index > max {
			max = col.index
		}
	}

	return max
}
