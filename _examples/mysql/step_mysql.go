package main

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/leeyaf/yflow"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	yflow.RegistStep("MysqlExecute", sMysqlExecute)
	yflow.RegistStep("MysqlExecuteParallel", sMysqlExecuteParallel)
}

// 执行 SQL
//
// 可以保证列的顺序与查询顺序一致，默认丢弃列名
//
// 基于 gorm.io/gorm 实现
func sMysqlExecute(s *yflow.Step, in *yflow.Matrix, out *yflow.Matrix) {
	uri := in.Get(0, 0) // 连接字符串
	sql := in.Get(1, 0) // SQL

	slog.Info("MysqlExecute", "sql", sql)

	executeSql(uri, sql, out)
}

func executeSql(uri string, sql string, out *yflow.Matrix) {
	db, err := gorm.Open(mysql.Open(uri))
	if err != nil {
		panic(err)
	}
	if sqldb, err := db.DB(); err != nil {
		panic(err)
	} else {
		sqldb.SetMaxIdleConns(1)
		defer sqldb.Close()
	}

	rows, err := db.Raw(sql).Rows()
	if err != nil {
		panic(err)
	}

	var columnNames []string
	for rows.Next() {
		if columnNames == nil {
			if names, err := rows.Columns(); err != nil {
				panic(err)
			} else {
				columnNames = names
			}
		}

		rowData := make(map[string]any)
		if err := db.ScanRows(rows, &rowData); err != nil {
			panic(err)
		}

		sortedRowData := make([]string, 0, len(columnNames))
		for _, name := range columnNames {
			sortedRowData = append(sortedRowData, fmt.Sprintf("%v", rowData[name]))
		}
		out.AppendRow(sortedRowData)
	}
}

// 并行向多个链接地址执行相同的 SQL，
// 并按行合并结果，不保证合并的顺序
//
// 适合滚服游戏：多张表的结构相同，需要跨表查询并合并结果
func sMysqlExecuteParallel(s *yflow.Step, in *yflow.Matrix, out *yflow.Matrix) {
	uriMatrixPath := in.Get(0, 0) // 存放连接字符串的矩阵名，每行一个链接地址
	sql := in.Get(1, 0)           // SQL

	uris := s.GetMatrix(uriMatrixPath).GetCol(0)

	var wg sync.WaitGroup
	wg.Add(len(uris))

	resultsChan := make(chan *yflow.Matrix, len(uris))
	for _, u := range uris {
		go func(uri string) {
			defer wg.Done()

			slog.Info("MysqlExecuteParallel", "uri", uri, "sql", sql)

			resultMatrix := yflow.NewMatrix()
			executeSql(uri, sql, resultMatrix)
			resultsChan <- resultMatrix
		}(u)
	}

	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		rows, _ := result.Shape()
		for row := range rows {
			out.AppendRow(result.GetRow(row))
		}
	}
}
