package yflow

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
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

func (c *column) copyTo(target *column) {
	target.data = c.slice()
}

func (c *column) slice() []string {
	newData := make([]string, len(c.data))
	copy(newData, c.data)
	return newData
}

func (c *column) append(data string) {
	c.data = append(c.data, data)
}

func (c *column) resize(newSize int, defaultValue string) {
	oldSize := len(c.data)

	if newSize < 0 {
		panic("negative size")
	} else if newSize > oldSize { // 扩大
		for range newSize - oldSize {
			c.data = append(c.data, defaultValue)
		}
	} else if newSize < oldSize { // 剪裁
		c.data = c.data[:newSize]
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

func (c *column) getData() []string {
	return c.data
}

func (c *column) setData(data []string) {
	c.data = data
}

// 基于列的二维矩阵
//
// 元素类型统一为 string
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
	// 在对矩阵进行操作时，需要确保所有 column 的长度一致
	rows = m.columns[0].len()
	return
}

func (m *Matrix) Clone() *Matrix {
	newColumns := make([]*column, 0, len(m.columns))
	for _, column := range m.columns {
		newColumns = append(newColumns, column.clone())
	}
	return &Matrix{
		columns: newColumns,
	}
}

func (m *Matrix) CopyTo(target *Matrix) {
	newColumns := make([]*column, 0, len(m.columns))
	for _, column := range m.columns {
		newColumns = append(newColumns, column.clone())
	}
	target.columns = newColumns
}

func (m *Matrix) Resize(rows int, cols int, defaultValue string) {
	if rows < 0 || cols < 0 {
		panic("negative rows or cols")
	}

	// 保证所有 column 的长度一致
	oldCols := 0
	for _, column := range m.columns {
		column.resize(rows, defaultValue)
		oldCols++
	}

	if cols > oldCols { // 扩大
		for range cols - oldCols {
			column := newColumn()
			column.resize(rows, defaultValue)
			m.columns = append(m.columns, column)
		}
	} else if cols < oldCols { // 剪裁
		m.columns = m.columns[:cols]
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

// 转置
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

// 根据矩阵的某列排序，该列的数据转成 float64 后排序
//
// asc 为 true 时升序，false 降序
func (m *Matrix) Sort(col int, asc bool) {
	if col < 0 || col >= len(m.columns) {
		panic("wrong col")
	}

	type dataWithRow struct {
		oldRow int
		data   float64
	}

	sortData := m.columns[col].slice()
	dwrs := make([]*dataWithRow, 0, len(sortData))
	for row, dataString := range sortData {
		dataFloat64, err := strconv.ParseFloat(dataString, 64)
		if err != nil {
			panic("parse float64 fail")
		}
		dwrs = append(dwrs, &dataWithRow{
			oldRow: row,
			data:   dataFloat64,
		})
	}

	sort.Slice(dwrs, func(i, j int) bool {
		if asc {
			return dwrs[i].data < dwrs[j].data
		} else {
			return dwrs[i].data > dwrs[j].data
		}
	})

	oldRowToNewRow := make(map[int]int, len(dwrs))
	for newRow, dwr := range dwrs {
		oldRowToNewRow[dwr.oldRow] = newRow
	}

	newColumns := make([]*column, 0, len(m.columns))
	for _, oldColumn := range m.columns {
		newColumn := newColumn()
		newColumn.resize(oldColumn.len(), "")
		for oldRow, cell := range oldColumn.getData() {
			newRow := oldRowToNewRow[oldRow]
			newColumn.set(newRow, cell)
		}
		newColumns = append(newColumns, newColumn)
	}

	m.columns = newColumns
}

// 插入一行
func (m *Matrix) InsertRow(row int, data []string) {
	if row < 0 {
		panic("negative row")
	}

	rows, cols := m.Shape()
	m.Resize(rows+1, cols, "")

	// 只移动插入行之后的行，每行的数据等于上一行的数据
	for i := rows; i > row; i-- {
		for _, column := range m.columns {
			column.set(i, column.get(i-1))
		}
	}

	// 插入该行
	for i, column := range m.columns {
		column.set(row, data[i])
	}
}

// 插入一列
func (m *Matrix) InsertCol(col int, data []string) {
	if col < 0 {
		panic("negative col")
	}

	rows, cols := m.Shape()
	m.Resize(rows, cols+1, "")

	// 只移动插入列之后的列，每列的数据等于前一列的数据
	for i := cols; i > col; i-- {
		m.columns[i] = m.columns[i-1]
	}

	// 插入该列
	m.columns[col].setData(data)
}

// 设置某行的数据
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

// 设置某列的数据
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

// 在末尾添加一行
func (m *Matrix) AppendRow(data []string) {
	rows, cols := m.Shape()
	m.Resize(rows, max(cols, len(data)), "")

	for i, col := range m.columns {
		col.append(data[i])
	}
}

// 在末尾添加一列
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

func (m *Matrix) GetInt(row int, col int) int {
	return int(m.GetInt64(row, col))
}

func (m *Matrix) GetInt32(row int, col int) int32 {
	return int32(m.GetInt64(row, col))
}

func (m *Matrix) GetInt64(row int, col int) int64 {
	cell := m.Get(row, col)
	i, err := strconv.ParseInt(cell, 10, 64)
	if err != nil {
		panic("parse int fail")
	}
	return i
}

func (m *Matrix) Set(row int, col int, value string) {
	rows, cols := m.Shape()

	if row < 0 || col < 0 || row >= rows || col >= cols {
		panic("wrong row or col")
	}

	m.columns[col].set(row, value)
}
