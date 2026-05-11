package tests

import (
	"fmt"
	"testing"

	"github.com/leeyaf/yflow"
)

func TestMatrix(t *testing.T) {
	m := yflow.NewMatrix()

	m.Resize(0, 0, "0")
	fmt.Println(m)

	m.Resize(11, 11, "-1")
	fmt.Println(m)

	m.Resize(12, 5, "good")
	fmt.Println(m)

	m.Resize(5, 12, "3.14")
	fmt.Println(m)

	m.Resize(12, 12, "6")
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
