package yflow

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/shlex"
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
	RegistStep("MatrixMerge", sMatrixMerge)

	RegistStep("Http", sHttp)
}

// 输出等于输入
//
// 可以用来测试表达式
func sTest(s *Step, in *Matrix, out *Matrix) {
	in.CopyTo(out)
}

// 从 YAML 创建矩阵
func sMatrixNew(s *Step, in *Matrix, out *Matrix) {
	rows, _ := in.Shape()
	for i := range rows {
		tokens, err := shlex.Split(in.Get(i, 0))
		if err != nil {
			panic(err)
		}
		out.AppendRow(tokens)
	}
}

// 创建一个空矩阵，并使用默认值填充
func sMatrixEmpty(s *Step, in *Matrix, out *Matrix) {
	rows := in.GetInt(0, 0)      // 行数
	cols := in.GetInt(1, 0)      // 列数
	defaultValue := in.Get(2, 0) // 默认值

	out.Resize(rows, cols, defaultValue)
}

// 矩阵转置
func sMatrixTranspose(s *Step, in *Matrix, out *Matrix) {
	matrixPath := in.Get(0, 0) // 要操作的矩阵名

	s.GetMatrix(matrixPath).CopyTo(out)
	out.Transpose()
}

// 根据矩阵的某列排序，该列的数据转成 float64 后排序
func sMatrixSort(s *Step, in *Matrix, out *Matrix) {
	matrixPath := in.Get(0, 0) // 要操作的矩阵名
	col := in.GetInt(1, 0)     // 列下标
	ascOrDesc := in.Get(2, 0)  // 升序或降序

	s.GetMatrix(matrixPath).CopyTo(out)
	out.Sort(col, ascOrDesc == "asc")
}

// 插入一行或一列到矩阵
func sMatrixInsert(s *Step, in *Matrix, out *Matrix) {
	matrixPath := in.Get(0, 0) // 要操作的矩阵名
	rowOrCol := in.Get(1, 0)   // 方向
	index := in.GetInt(2, 0)   // 下标
	values := in.Get(3, 0)     // 数据

	tokens, err := shlex.Split(values)
	if err != nil {
		panic(err)
	}

	s.GetMatrix(matrixPath).CopyTo(out)

	if rowOrCol == "row" {
		out.InsertRow(index, tokens)
	} else if rowOrCol == "col" {
		out.InsertCol(index, tokens)
	} else {
		panic("wrong row or col")
	}
}

// 筛选矩阵的多行或多列
func sMatrixSelect(s *Step, in *Matrix, out *Matrix) {
	matrixPath := in.Get(0, 0) // 原矩阵名
	rowOrCol := in.Get(1, 0)   // 方向
	indexs := in.Get(2, 0)     // 下标数组

	tokens, err := shlex.Split(indexs)
	if err != nil {
		panic(err)
	}

	sourceMatrix := s.GetMatrix(matrixPath)
	for _, token := range tokens {
		index, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			panic(err)
		}
		if rowOrCol == "row" {
			rowData := sourceMatrix.GetRow(int(index))
			out.AppendRow(rowData)
		} else if rowOrCol == "col" {
			colData := sourceMatrix.GetCol(int(index))
			out.AppendCol(colData)
		} else {
			panic("wrong name, row or col")
		}
	}
}

// 垂直查找并替换矩阵的一列
func sMatrixVLookUp(s *Step, in *Matrix, out *Matrix) {
	sourceMatrixPath := in.Get(0, 0)   // 原矩阵名
	sourceCol := in.GetInt(1, 0)       // 原矩阵的列下标
	lookUpMatrixPath := in.Get(2, 0)   // 目标矩阵名
	lookUpCol := in.GetInt(3, 0)       // 目标矩阵查找列的下标
	lookUpReturnCol := in.GetInt(4, 0) // 目标矩阵返回列的下标

	sMat := s.GetMatrix(sourceMatrixPath)
	lMat := s.GetMatrix(lookUpMatrixPath)

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
func sMatrixApply(s *Step, in *Matrix, out *Matrix) {
	sourceMatrixPath := in.Get(0, 0) // 要操作的矩阵名
	rowOrCol := in.Get(1, 0)         // 方向
	index := in.GetInt(2, 0)         // 下标
	varName := in.Get(3, 0)          // 元素的变量名
	expression := in.Get(4, 0)       // 表达式

	s.GetMatrix(sourceMatrixPath).CopyTo(out)

	if rowOrCol == "row" {
		newSlices := make([]string, 0)
		for _, value := range out.GetRow(index) {
			s.GengineAddContext(varName, value)
			result := s.GengineExecute(expression)
			newSlices = append(newSlices, result)
		}
		out.SetRow(index, newSlices)
	} else if rowOrCol == "col" {
		newSlices := make([]string, 0)
		for _, value := range out.GetCol(index) {
			s.GengineAddContext(varName, value)
			result := s.GengineExecute(expression)
			newSlices = append(newSlices, result)
		}
		out.SetCol(index, newSlices)
	} else {
		panic("wrong row or col")
	}
}

// 按行或按列合并两个矩阵
func sMatrixMerge(s *Step, in *Matrix, out *Matrix) {
	matrixPath1 := in.Get(0, 0) // 源矩阵1
	matrixPath2 := in.Get(1, 0) // 源矩阵2
	rowOrCol := in.Get(2, 0)    // 方向

	mergeMatrix := s.GetMatrix(matrixPath2)
	rows, cols := mergeMatrix.Shape()

	s.GetMatrix(matrixPath1).CopyTo(out)

	switch rowOrCol {
	case "row":
		for row := range rows {
			out.AppendRow(mergeMatrix.GetRow(row))
		}
	case "col":
		for col := range cols {
			out.AppendCol(mergeMatrix.GetCol(col))
		}
	default:
		panic("wrong row or col")
	}
}

func sHttp(s *Step, in *Matrix, out *Matrix) {
	method := in.Get(0, 0) // 留空代表 GET，可选：OPTIONS GET HEAD POST PUT DELETE TRACE CONNECT
	uri := in.Get(1, 0)    // 地址
	header := in.Get(2, 0) // 使用=连接键值 空格分隔多组 Header
	body := in.Get(3, 0)   // 字符串

	req, err := http.NewRequest(method, uri, bytes.NewBufferString(body))
	if err != nil {
		panic(err)
	}

	// 处理 header
	if len(header) > 0 {
		tokens, err := shlex.Split(header)
		if err != nil {
			panic(err)
		}
		for _, token := range tokens {
			parts := strings.Split(token, "=")
			req.Header.Set(parts[0], parts[1])
		}
	}

	cli := &http.Client{
		Timeout: StepMaxLifetime,
	}
	rsp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	rspData, err := io.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}

	out.Set(0, 0, strconv.Itoa(rsp.StatusCode))
	out.Set(1, 0, string(rspData))
}
