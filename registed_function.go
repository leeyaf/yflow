package yflow

import (
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
}

func fCell(step *Step) any {
	return func(matrixName string, row int, col int) string {
		return step.GetMatrix(matrixName).Get(row, col)
	}
}

func fRow(step *Step) any {
	return func(matrixName string, row int) []string {
		return step.GetMatrix(matrixName).GetRow(row)
	}
}

func fCol(step *Step) any {
	return func(matrixName string, col int) []string {
		return step.GetMatrix(matrixName).GetCol(col)
	}
}

func fParseInt(val string) int64 {
	if intVal, err := strconv.ParseInt(val, 10, 64); err != nil {
		panic(err)
	} else {
		return intVal
	}
}

func fParseFloat(val string) float64 {
	if floatVal, err := strconv.ParseFloat(val, 64); err != nil {
		panic(err)
	} else {
		return floatVal
	}
}

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

func fSplit(s string, sep rune) []string {
	return strings.Split(s, string(sep))
}

// time.String returns the time formatted using the format string
//
//	"2006-01-02 15:04:05.999999999 -0700 MST"
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
