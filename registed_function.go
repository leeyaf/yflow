package yflow

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegistFunctionWithStep("Cell", fCell)
	RegistFunctionWithStep("Row", fRow)
	RegistFunctionWithStep("Rows", fRows)
	RegistFunctionWithStep("Col", fCol)
	RegistFunctionWithStep("Cols", fCols)
	RegistFunction("ParseInt", fParseInt)
	RegistFunction("ParseFloat", fParseFloat)
	RegistFunction("Join", fJoin)
	RegistFunction("Split", fSplit)
	RegistFunction("TimeConvert", fTimeConvert)
	RegistFunction("Sprintf", fSprintf)
	RegistFunction("Arange", fArange)
}

// 从矩阵中读取值
func fCell(step *Step) any {
	return func(matrixPath string, row int, col int) string {
		return step.GetMatrix(matrixPath).Get(row, col)
	}
}

// 从矩阵中读取一行
func fRow(step *Step) any {
	return func(matrixPath string, row int) []string {
		return step.GetMatrix(matrixPath).GetRow(row)
	}
}

// 获取矩阵的行数
func fRows(step *Step) any {
	return func(matrixPath string) int {
		rows, _ := step.GetMatrix(matrixPath).Shape()
		return rows
	}
}

// 从矩阵中读取一列
func fCol(step *Step) any {
	return func(matrixPath string, col int) []string {
		return step.GetMatrix(matrixPath).GetCol(col)
	}
}

// 获取矩阵的列数
func fCols(step *Step) any {
	return func(matrixPath string) int {
		_, cols := step.GetMatrix(matrixPath).Shape()
		return cols
	}
}

// 字符串转 int64
func fParseInt(val string) int64 {
	if intVal, err := strconv.ParseInt(val, 10, 64); err != nil {
		panic(err)
	} else {
		return intVal
	}
}

// 字符串转 float64
func fParseFloat(val string) float64 {
	if floatVal, err := strconv.ParseFloat(val, 64); err != nil {
		panic(err)
	} else {
		return floatVal
	}
}

// 切片转字符串
//
// sep 用来分隔元素
//
// wrap 用来包裹元素
//
// 注意：无法在 YAML 中使用 " 作为 sep 或 wrap
func fJoin(datas []string, sep string, wrap string) string {
	wrapedDatas := make([]string, 0, len(datas))
	for _, data := range datas {
		wrapedDatas = append(wrapedDatas, wrap+data+wrap)
	}
	return strings.Join(wrapedDatas, sep)
}

// 分隔字符串
//
// sep 分隔符
func fSplit(s string, sep string) []string {
	return strings.Split(s, sep)
}

// 转换时间字符串的格式
//
// oldLayout 原模版，newLayout 新模版
//
// go 使用 "2006-01-02 15:04:05.999999999 -0700 MST" 转换 time
func fTimeConvert(s string, oldLayout string, newLayout string) string {
	if oldLayout == "" {
		oldLayout = "2006-01-02 15:04:05.999999999 -0700 MST"
	}
	t, err := time.Parse(oldLayout, s)
	if err != nil {
		panic(err)
	}
	return t.Format(newLayout)
}

// fmt.Sprintf()
func fSprintf(s string, a ...any) string {
	return fmt.Sprintf(s, a...)
}

// 生成从 start 到 stop（不含）的 int 切片，步长 step
//
// 类似 NumPy 的 arange
func fArange(start, stop, step int) []string {
	if step == 0 {
		panic("step cannot be zero")
	}

	length := 0
	if (step > 0 && start < stop) || (step < 0 && start > stop) {
		length = max(0, (stop-start+step-1)/step)
	}

	result := make([]string, 0, length)
	for i := range length {
		result = append(result, strconv.Itoa(start+i*step))
	}
	return result
}
