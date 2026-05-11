package yflow

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	printMax = 11
)

type column struct {
	data []string
}

func newColumn() *column {
	return &column{
		data: make([]string, 0),
	}
}

func (c *column) clone() *column {
	return &column{
		data: c.slice(),
	}
}

func (c *column) slice() []string {
	newData := make([]string, len(c.data))
	copy(newData, c.data)
	return newData
}

func (c *column) append(data string) {
	c.data = append(c.data, data)
}

func (c *column) resize(size int, defaultValue string) {
	if size < 0 {
		panic("negative size")
	} else if len(c.data) < size {
		// 扩大
		for range size - len(c.data) {
			c.data = append(c.data, defaultValue)
		}
	} else if len(c.data) > size {
		// 剪裁
		c.data = c.data[:size]
	} else {
		return
	}
}

func (c *column) sum() string {
	bigTotal := new(big.Float)
	for _, value := range c.data {
		bigValue := new(big.Float)
		bigValue.SetString(value)
		bigTotal.Add(bigTotal, bigValue)
	}
	return bigTotal.String()
}

func (c *column) avg() string {
	bigTotal := new(big.Float)
	bigTotal.SetString(c.sum())

	bigCount := new(big.Float)
	bigCount.SetInt64((int64)(len(c.data)))

	bigAvg := new(big.Float)
	bigAvg.Quo(bigTotal, bigCount)
	return bigAvg.String()
}

func (c *column) mathOperation(ope func(d *big.Float, v *big.Float), v string) {
	bigValue := new(big.Float)
	bigValue.SetString(v)

	for i, d := range c.data {
		bigD := new(big.Float)
		bigD.SetString(d)
		ope(bigD, bigValue)
		c.data[i] = bigD.String()
	}
}

func (c *column) add(value string) {
	c.mathOperation(func(d, v *big.Float) {
		d.Add(d, v)
	}, value)
}

func (c *column) sub(value string) {
	c.mathOperation(func(d, v *big.Float) {
		d.Sub(d, v)
	}, value)
}

func (c *column) mul(value string) {
	c.mathOperation(func(d, v *big.Float) {
		d.Mul(d, v)
	}, value)
}

func (c *column) quo(value string) {
	c.mathOperation(func(d, v *big.Float) {
		d.Quo(d, v)
	}, value)
}

func (c *column) len() int {
	return len(c.data)
}

func (c *column) get(row int) string {
	return c.data[row]
}

func (c *column) set(row int, value string) {
	c.data[row] = value
}

func (c *column) setData(data []string) {
	c.data = data
}

type Matrix struct {
	columns []*column
}

func NewMatrix() *Matrix {
	return &Matrix{
		columns: make([]*column, 0),
	}
}

func (m *Matrix) Shape() (rows, cols int) {
	if len(m.columns) == 0 {
		return
	}
	cols = len(m.columns)
	rows = m.columns[0].len() // 确保所有的 col 长度一致
	return
}

func (m *Matrix) Clone() *Matrix {
	newColumns := make([]*column, 0, len(m.columns))
	for _, col := range m.columns {
		newColumns = append(newColumns, col.clone())
	}
	return &Matrix{
		columns: newColumns,
	}
}

func (m *Matrix) SetColumns(columns []*column) {
	m.columns = columns
}

func (m *Matrix) GetColumns() []*column {
	return m.columns
}

func (m *Matrix) Resize(rows int, cols int, defaultValue string) {
	if rows < 0 || cols < 0 {
		panic("negative rows or cols")
	}

	_, oldCols := m.Shape()
	if cols > oldCols {
		// 扩大
		for range cols - oldCols {
			column := newColumn()
			column.resize(cols, defaultValue)
			m.columns = append(m.columns, column)
		}
	} else if cols < oldCols {
		// 剪裁
		m.columns = m.columns[:cols-1]
	}

	for _, col := range m.columns {
		col.resize(rows, defaultValue)
	}
}

func (m *Matrix) String() string {
	rows, cols := m.Shape()
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf("(%v, %v)\n", rows, cols))
	if rows <= printMax {
		for row := range rows {
			m.rowString(b, row)
		}
	} else {
		for row := range (printMax - 1) / 2 {
			m.rowString(b, row)
		}
		b.WriteString("..\n")
		for row := rows - (printMax-1)/2; row < rows; row++ {
			m.rowString(b, row)
		}
	}
	return b.String()
}

func (m *Matrix) rowString(b *strings.Builder, row int) {
	rows, cols := m.Shape()

	breaker := ""
	if row+1 < rows {
		breaker = "\n"
	}
	if cols <= printMax {
		m.subRowString(b, row, 0, cols, breaker)
	} else {
		m.subRowString(b, row, 0, (printMax-1)/2, ", .., ")
		m.subRowString(b, row, cols-(printMax-1)/2, cols, breaker)
	}
}

// (colFrom, colTo]
func (m *Matrix) subRowString(b *strings.Builder, row int, colFrom, colTo int, breaker string) {
	for i := colFrom; i < colTo; i++ {
		b.WriteString(m.columns[i].get(row))
		if i+1 < colTo {
			b.WriteString(", ")
		} else {
			b.WriteString(breaker)
		}
	}
}

func (m *Matrix) Transpose() {
	rows, cols := m.Shape()

	if c := rows + cols; c < 3 {
		return // (0, 0) or (1, 1)
	}

	oldColumns := m.columns

	newColumns := make([]*column, 0, rows)
	for range rows {
		column := newColumn()
		column.resize(cols, "")
		newColumns = append(newColumns, column)
	}

	for row := range cols {
		for col := range rows {
			newColumns[col].set(row, oldColumns[row].get(col))
		}
	}
	m.columns = newColumns
}

func (m *Matrix) InsertRow(row int, data []string) {
	if row < 0 {
		panic("negative row")
	}

	rows, cols := m.Shape()
	if len(data) != cols {
		panic("wrong shape")
	}
	m.Resize(rows+1, cols, "")

	for i := rows; i > row; i-- {
		for _, column := range m.columns {
			column.set(i, column.get(i-1))
		}
	}
	for i, column := range m.columns {
		column.set(row, data[i])
	}
}

func (m *Matrix) InsertCol(col int, data []string) {
	if col < 0 {
		panic("negative col")
	}

	rows, cols := m.Shape()
	if len(data) != rows {
		panic("wrong shape")
	}
	m.Resize(rows, cols+1, "")

	for i := cols; i > col; i-- {
		m.columns[i] = m.columns[i-1]
	}
	m.columns[col].setData(data)
}

func (m *Matrix) SetRow(row int, data []string) {
	if row < 0 {
		panic("negative row")
	}

	rows, cols := m.Shape()
	m.Resize(max(rows, row+1), max(cols, len(data)), "")

	for j := range len(data) {
		m.columns[j].set(row, data[j])
	}
}

func (m *Matrix) SetCol(col int, data []string) {
	if col < 0 {
		panic("negative row")
	}

	rows, cols := m.Shape()
	m.Resize(max(rows, len(data)), max(cols, col+1), "")

	for i := range data {
		m.columns[col].set(i, data[i])
	}
}

func (m *Matrix) AppendRow(data []string) {
	rows, cols := m.Shape()
	m.Resize(rows, max(cols, len(data)), "")

	for i, col := range m.columns {
		col.append(data[i])
	}
}

func (m *Matrix) AppendCol(data []string) {
	rows, cols := m.Shape()
	m.Resize(max(rows, len(data)), cols, "")

	column := newColumn()
	column.setData(data)
	m.columns = append(m.columns, column)
}

func (m *Matrix) GetRow(row int) []string {
	rows, cols := m.Shape()

	if row < 0 || row >= rows {
		panic("wrong row")
	}

	rowData := make([]string, 0, cols)
	for _, col := range m.columns {
		rowData = append(rowData, col.get(row))
	}
	return rowData
}

func (m *Matrix) GetCol(col int) []string {
	_, cols := m.Shape()

	if col < 0 || col >= cols {
		panic("wrong col")
	}

	return m.columns[col].slice()
}

func (m *Matrix) Get(row int, col int) string {
	rows, cols := m.Shape()

	if row < 0 || col < 0 || row >= rows || col >= cols {
		panic("wrong row or col")
	}

	return m.columns[col].get(row)
}

func (m *Matrix) Set(row int, col int, value string) {
	rows, cols := m.Shape()

	if row < 0 || col < 0 || row >= rows || col >= cols {
		panic("wrong row or col")
	}

	m.columns[col].set(row, value)
}
