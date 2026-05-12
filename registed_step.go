package yflow

import (
	"strconv"
	"strings"
)

func init() {
	RegistStep("Test", sTest)

	RegistStep("MatrixNew", sMatrixNew)
	RegistStep("MatrixEmpty", sMatrixEmpty)
	RegistStep("MatrixTranspose", sMatrixTranspose)
	RegistStep("MatrixSort", sMatrixSort)
	RegistStep("MatrixInsert", sMatrixInsert)
	RegistStep("MatrixSelect", sMatrixSelect)
	RegistStep("MatrixVLookUp", sMatrixVLookUp)
	RegistStep("MatrixApply", sMatrixApply)
}

// 输出等于输入
//
// 可以用来测试表达式
func sTest(s *Step, in *Matrix, out *Matrix) {
	in.CopyTo(out)
}

// 从 YAML 定义创建矩阵
func sMatrixNew(s *Step, in *Matrix, out *Matrix) {
	sep := in.Get(0, 0)      // 行内元素的字符串分隔符
	rows := in.GetCol(0)[1:] // 多行数据

	for _, row := range rows {
		cells := strings.Split(row, sep)
		trimedCells := make([]string, 0, len(cells))
		for _, c := range cells {
			trimedCells = append(trimedCells, strings.TrimSpace(c))
		}
		out.AppendRow(trimedCells)
	}
}

// 创建一个空矩阵，并使用默认值填充
func sMatrixEmpty(s *Step, in *Matrix, out *Matrix) {
	rows := in.GetInt(0, 0)      // 行数
	cols := in.GetInt(1, 0)      // 列数
	defaultValue := in.Get(2, 0) // 默认值

	out.Resize(rows, cols, defaultValue)
}

// 矩阵转置，行变成列，列变成行
func sMatrixTranspose(s *Step, in *Matrix, out *Matrix) {
	matrixName := in.Get(0, 0) // 要操作的矩阵名

	s.GetMatrix(matrixName).CopyTo(out)
	out.Transpose()
}

// 根据矩阵的某列排序，该列的数据转成 float64 后排序
func sMatrixSort(s *Step, in *Matrix, out *Matrix) {
	matrixName := in.Get(0, 0) // 要操作的矩阵名
	col := in.GetInt(1, 0)     // 列下标
	ascOrDesc := in.Get(2, 0)  // 升序或降序

	s.GetMatrix(matrixName).CopyTo(out)
	out.Sort(col, ascOrDesc == "asc")
}

// 插入一行或一列到矩阵
func sMatrixInsert(s *Step, in *Matrix, out *Matrix) {
	matrixName := in.Get(0, 0) // 要操作的矩阵名
	rowOrCol := in.Get(1, 0)   // 方向
	index := in.GetInt(2, 0)   // 下标
	values := in.GetCol(0)[3:] // 数据

	s.GetMatrix(matrixName).CopyTo(out)

	if rowOrCol == "row" {
		out.InsertRow(index, values)
	} else if rowOrCol == "col" {
		out.InsertCol(index, values)
	} else {
		panic("wrong row or col")
	}
}

// 筛选矩阵的多行或多列
func sMatrixSelect(s *Step, in *Matrix, out *Matrix) {
	sourceMatrixName := in.Get(0, 0) // 原矩阵名
	rowOrCol := in.Get(1, 0)         // 方向
	indexs := in.GetCol(0)[2:]       // 下标数组

	sourceMatrix := s.GetMatrix(sourceMatrixName)
	for _, index := range indexs {
		i, err := strconv.ParseInt(index, 10, 64)
		if err != nil {
			panic(err)
		}
		if rowOrCol == "row" {
			rowData := sourceMatrix.GetRow(int(i))
			out.AppendRow(rowData)
		} else if rowOrCol == "col" {
			colData := sourceMatrix.GetCol(int(i))
			out.AppendCol(colData)
		} else {
			panic("wrong name, row or col")
		}
	}
}

// 垂直查找并替换矩阵的一列
func sMatrixVLookUp(s *Step, in *Matrix, out *Matrix) {
	sourceMatrix := in.Get(0, 0)       // 原矩阵名
	sourceCol := in.GetInt(1, 0)       // 原矩阵的列下标
	lookUpMatrix := in.Get(2, 0)       // 目标矩阵名
	lookUpCol := in.GetInt(3, 0)       // 目标矩阵查找列的下标
	lookUpReturnCol := in.GetInt(4, 0) // 目标矩阵返回列的下标

	sMat := s.GetMatrix(sourceMatrix)
	lMat := s.GetMatrix(lookUpMatrix)

	sColumn := sMat.GetCol(sourceCol)
	lColumn := lMat.GetCol(lookUpCol)
	rColumn := lMat.GetCol(lookUpReturnCol)

	find := func(source string, lookUpCol []string, returnCol []string) string {
		for i := range lookUpCol {
			l := lookUpCol[i]

			if l == source {
				return returnCol[i]
			}
		}
		return source
	}

	newColunm := make([]string, 0, len(sColumn))
	for _, origin := range sColumn {
		replace := find(origin, lColumn, rColumn)
		newColunm = append(newColunm, replace)
	}

	sMat.CopyTo(out)

	out.SetCol(sourceCol, newColunm)
}

// 对矩阵的行或列应用表达式
//
// 注意：表达式不要使用 (( )) 包裹
// 包裹会按照 step.in 的逻辑统一处理
//
// value 代表元素的值
func sMatrixApply(s *Step, in *Matrix, out *Matrix) {
	sourceMatrixName := in.Get(0, 0) // 要操作的矩阵名
	rowOrCol := in.Get(1, 0)         // 方向
	index := in.GetInt(2, 0)         // 下标
	expression := in.Get(3, 0)       // 表达式

	s.GetMatrix(sourceMatrixName).CopyTo(out)

	if rowOrCol == "row" {
		newSlices := make([]string, 0)
		for _, value := range out.GetRow(index) {
			result := s.Apply(value, "value", expression)
			newSlices = append(newSlices, result)
		}
		out.SetRow(index, newSlices)
	} else if rowOrCol == "col" {
		newSlices := make([]string, 0)
		for _, value := range out.GetCol(index) {
			result := s.Apply(value, "value", expression)
			newSlices = append(newSlices, result)
		}
		out.SetCol(index, newSlices)
	} else {
		panic("wrong row or col")
	}
}
