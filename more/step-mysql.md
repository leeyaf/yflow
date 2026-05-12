# StepMysql

## MysqlExecute

基于 [gorm.io/gorm](https://github.com/go-gorm/gorm) 实现对 MySQL 的操作

可以保证列的顺序与查询顺序一致，默认丢弃列名

```` go
import (
	"fmt"
	"log/slog"

	"github.com/leeyaf/yflow"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	yflow.RegistStep(
		"MysqlExecute",
		[]string{"Uri", "Sql"},
		sMysqlExecute,
	)
}

// 执行 MySQL 查询
func sMysqlExecute(s *yflow.Step, in *yflow.Matrix, out *yflow.Matrix) {
	uri := in.Get(0, 0) // 连接字符串
	sql := in.Get(1, 0) // SQL

	slog.Info("MysqlExecute", "sql", sql)

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
````