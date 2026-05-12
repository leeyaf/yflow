package tests

import (
	"fmt"
	"testing"

	"github.com/leeyaf/yflow"
)

func TestMatrixResize(t *testing.T) {
	m := yflow.NewMatrix()

	m.Resize(6, 6, "0")
	fmt.Println(m)

	m.Resize(3, 3, "-1")
	fmt.Println(m)

	m.Resize(5, 5, "a")
	fmt.Println(m)
}

func TestMatrixTranspose(t *testing.T) {
	m := yflow.NewMatrix()

	m.Resize(3, 4, "0")
	m.SetRow(0, []string{"0", "4", "5", "6"})
	m.SetCol(0, []string{"1", "2", "3"})
	fmt.Println(m)

	m.Transpose()
	fmt.Println(m)

	m.Transpose()
	fmt.Println(m)
}

func TestMatrixSum(t *testing.T) {
	m := yflow.NewMatrix()

	m.Resize(3, 1, "-0.5")
	fmt.Println(m)

	m.Resize(3, 1, "abc")
	fmt.Println(m)

	m.Resize(3, 1, "1.1")
	m.Set(1, 0, "1.2")
	m.Set(2, 0, "1.34")
	fmt.Println(m)
}

func TestMatrixSort(t *testing.T) {
	matrix := yflow.NewMatrix()
	matrix.AppendRow([]string{"a1", "-0.5", "c1"})
	matrix.AppendRow([]string{"a2", "0", "c2"})
	matrix.AppendRow([]string{"a3", "-10", "c3"})
	matrix.AppendRow([]string{"a4", "20", "c4"})
	matrix.AppendRow([]string{"a5", "1.3", "c5"})
	matrix.Sort(1, true)
	fmt.Println(matrix)
	matrix.Sort(1, false)
	fmt.Println(matrix)
}
