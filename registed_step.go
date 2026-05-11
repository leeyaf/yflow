package yflow

import (
	"strconv"
	"strings"
)

func init() {
	RegistStep(
		"Test",
		[]string{"Param1", "Param2", ".."},
		sTest,
	)
	RegistStep(
		"MatrixNew",
		[]string{"Sep", "Row1", "Row2", ".."},
		sMatrixNew,
	)
	RegistStep(
		"MatrixTranspose",
		[]string{"SourceMatrixName"},
		sMatrixTranspose,
	)
	RegistStep(
		"MatrixInsert",
		[]string{"SourceMatrixName", "RowOrCol", "Index", "Value0", "Value1", ".."},
		sMatrixInsert,
	)
	RegistStep(
		"MatrixSelect",
		[]string{"SourceMatrixName", "RowOrCol", "Index0", "Index1", ".."},
		sMatrixSelect,
	)
	RegistStep(
		"MatrixVLookUp",
		[]string{"SourceMatrixName", "SourceCol", "LoopUpMatrixName", "LookUpCol", "ReturnCol"},
		sMatrixVLookUp,
	)
	RegistStep(
		"MatrixApply",
		[]string{"SourceMatrixName", "RowOrCol", "Index", "Expression"},
		sMatrixApply,
	)
}

func sTest(s *Step, in *Matrix, out *Matrix) {
	out.SetColumns(in.GetColumns())
}

func sMatrixNew(s *Step, in *Matrix, out *Matrix) {
	sep := in.Get(0, 0)
	rows := in.GetCol(0)[1:]

	for _, row := range rows {
		cells := strings.Split(row, sep)
		trimedCells := make([]string, 0, len(cells))
		for _, c := range cells {
			trimedCells = append(trimedCells, strings.TrimSpace(c))
		}
		out.AppendRow(trimedCells)
	}
}

func sMatrixTranspose(s *Step, in *Matrix, out *Matrix) {
	matrixName := in.Get(0, 0)

	sourceMatrix := s.GetMatrix(matrixName).Clone()
	out.SetColumns(sourceMatrix.GetColumns())
	out.Transpose()
}

func sMatrixInsert(s *Step, in *Matrix, out *Matrix) {
	col0 := in.GetCol(0)
	matrixName := col0[0]
	rowOrCol := col0[1]
	index := col0[2]
	values := col0[3:]

	intIndex, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		panic(err)
	}

	sourceMatrix := s.GetMatrix(matrixName).Clone()
	out.SetColumns(sourceMatrix.GetColumns())
	if rowOrCol == "row" {
		out.InsertRow(int(intIndex), values)
	} else if rowOrCol == "col" {
		out.InsertCol(int(intIndex), values)
	} else {
		panic("wrong row or col")
	}
}

func sMatrixSelect(s *Step, in *Matrix, out *Matrix) {
	col0 := in.GetCol(0)
	sourceMatrixName := col0[0]
	rowOrCol := col0[1]
	indexs := col0[2:]

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

func sMatrixVLookUp(s *Step, in *Matrix, out *Matrix) {
	sourceMatrix := in.Get(0, 0)
	sourceCol := in.Get(1, 0)
	lookUpMatrix := in.Get(2, 0)
	lookUpCol := in.Get(3, 0)
	lookUpReturnCol := in.Get(4, 0)

	sCol, err := strconv.ParseInt(sourceCol, 10, 64)
	if err != nil {
		panic(err)
	}

	lCol, err := strconv.ParseInt(lookUpCol, 10, 64)
	if err != nil {
		panic(err)
	}

	rCol, err := strconv.ParseInt(lookUpReturnCol, 10, 64)
	if err != nil {
		panic(err)
	}

	sMat := s.GetMatrix(sourceMatrix)
	lMat := s.GetMatrix(lookUpMatrix)

	sColumn := sMat.GetCol(int(sCol))
	lColumn := lMat.GetCol(int(lCol))
	rColumn := lMat.GetCol(int(rCol))

	newColunm := make([]string, 0, len(sColumn))
	for _, origin := range sColumn {
		replace := sMatrixVLookUpLook(origin, lColumn, rColumn)
		newColunm = append(newColunm, replace)
	}

	out.SetColumns(sMat.Clone().GetColumns())
	out.SetCol(int(sCol), newColunm)
}

func sMatrixVLookUpLook(source string, lookUpCol []string, returnCol []string) string {
	for i := range lookUpCol {
		l := lookUpCol[i]

		if l == source {
			return returnCol[i]
		}
	}
	return source
}

func sMatrixApply(s *Step, in *Matrix, out *Matrix) {
	sourceMatrixName := in.Get(0, 0)
	rowOrCol := in.Get(1, 0)
	index := in.Get(2, 0)
	expression := in.Get(3, 0)

	sMat := s.GetMatrix(sourceMatrixName)
	intIndex, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		panic(err)
	}

	out.SetColumns(sMat.Clone().GetColumns())

	if rowOrCol == "row" {
		newSlices := make([]string, 0)
		for _, value := range out.GetRow(int(intIndex)) {
			result := s.Apply(value, expression)
			newSlices = append(newSlices, result)
		}
		out.SetRow(int(intIndex), newSlices)
	} else if rowOrCol == "col" {
		newSlices := make([]string, 0)
		for _, value := range out.GetCol(int(intIndex)) {
			result := s.Apply(value, expression)
			newSlices = append(newSlices, result)
		}
		out.SetCol(int(intIndex), newSlices)
	} else {
		panic("wrong row or col")
	}
}
