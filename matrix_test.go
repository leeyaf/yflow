package yflow_test

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

func TestMatrixInsert(t *testing.T) {
	matrix := yflow.NewMatrix()
	matrix.AppendRow([]string{"a1", "b1", "c1"})
	matrix.AppendRow([]string{"a2", "b2", "c2"})
	matrix.AppendRow([]string{"a3", "b3", "c3"})

	fmt.Println("============")
	fmt.Println(matrix)
	rowMat := matrix.Clone()
	rowMat.InsertRow(1, []string{"short1", "short2"})
	fmt.Println(rowMat)
	rowMat.InsertRow(1, []string{"long1", "long1", "long3", "long4"})
	fmt.Println(rowMat)

	fmt.Println("============")
	fmt.Println(matrix)
	colMat := matrix.Clone()
	colMat.InsertCol(1, []string{"short1", "short2"})
	fmt.Println(colMat)
	colMat.InsertCol(1, []string{"long1", "long1", "long3", "long4"})
	fmt.Println(colMat)
}

func TestMatrixSet(t *testing.T) {
	matrix := yflow.NewMatrix()
	matrix.AppendRow([]string{"a1", "b1", "c1"})
	matrix.AppendRow([]string{"a2", "b2", "c2"})
	matrix.AppendRow([]string{"a3", "b3", "c3"})

	fmt.Println("============")
	fmt.Println(matrix)
	rowMat := matrix.Clone()
	rowMat.SetRow(0, []string{"short1", "short2"})
	fmt.Println(rowMat)
	rowMat.SetRow(0, []string{"long1", "long1", "long3", "long4"})
	fmt.Println(rowMat)
	rowMat.SetRow(4, []string{"skip1"})
	fmt.Println(rowMat)

	fmt.Println("============")
	fmt.Println(matrix)
	colMat := matrix.Clone()
	colMat.SetCol(0, []string{"short1", "short2"})
	fmt.Println(colMat)
	colMat.SetCol(0, []string{"long1", "long1", "long3", "long4"})
	fmt.Println(colMat)
	colMat.SetCol(4, []string{"skip1"})
	fmt.Println(colMat)
}

func TestMartixAppend(t *testing.T) {
	matrix := yflow.NewMatrix()
	matrix.AppendRow([]string{"a1", "b1", "c1"})
	matrix.AppendRow([]string{"a2", "b2", "c2"})
	matrix.AppendRow([]string{"a3", "b3", "c3"})

	fmt.Println("============")
	fmt.Println(matrix)
	rowMat := matrix.Clone()
	rowMat.AppendRow([]string{"short1", "short2"})
	fmt.Println(rowMat)
	rowMat.AppendRow([]string{"long1", "long1", "long3", "long4"})
	fmt.Println(rowMat)
	rowMat.AppendRow([]string{"one1"})
	fmt.Println(rowMat)

	fmt.Println("============")
	fmt.Println(matrix)
	colMat := matrix.Clone()
	colMat.AppendCol([]string{"short1", "short2"})
	fmt.Println(colMat)
	colMat.AppendCol([]string{"long1", "long1", "long3", "long4"})
	fmt.Println(colMat)
	colMat.AppendCol([]string{"one1"})
	fmt.Println(colMat)
}
