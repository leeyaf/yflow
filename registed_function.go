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
	RegistFunctionWithStep("Col", fCol)
	RegistFunction("ParseInt", fParseInt)
	RegistFunction("ParseFloat", fParseFloat)
	RegistFunction("Join", fJoin)
	RegistFunction("Split", fSplit)
	RegistFunction("TimeConvert", fTimeConvert)
	RegistFunction("Sprintf", fSprintf)
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

// 从矩阵中读取一列
func fCol(step *Step) any {
	return func(matrixPath string, col int) []string {
		return step.GetMatrix(matrixPath).GetCol(col)
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
// 注意：这里使用 rune 类型，默认值 0
//
// 常用符号与 ASCII 十进制的对应关系如下：
//
// “ 34
//
// ' 39
//
// , 44
//
// - 45
//
// . 46
//
// : 58
//
// ; 59
//
// ` 96
func fJoin(datas []string, asciiSep rune, asciiWrap rune) string {
	sep := string(asciiSep)
	if asciiWrap == 0 {
		return strings.Join(datas, sep)
	} else {
		wrap := string(asciiWrap)
		wrapedDatas := make([]string, 0, len(datas))
		for _, data := range datas {
			wrapedDatas = append(wrapedDatas, wrap+data+wrap)
		}
		return strings.Join(wrapedDatas, sep)
	}
}

// 分隔字符串
//
// sep 分隔符 rune 类型
func fSplit(s string, sep rune) []string {
	return strings.Split(s, string(sep))
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
