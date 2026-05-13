# 一个比较复杂的示例

## 准备 SQL

```` sql
CREATE DATABASE yflow;

USE yflow;

CREATE TABLE hero (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    player_id INT NOT NULL,
    config_id INT NOT NULL,
    level INT NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_player_id (player_id),
    INDEX idx_config_id (config_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO hero (player_id, config_id, level, created_time) VALUES
(1, 1, 1, '2026-01-15 10:30:00'),
(1, 2, 5, '2026-01-16 14:20:00'),
(2, 1, 3, '2026-01-17 09:15:00'),
(2, 3, 8, '2026-01-18 16:45:00'),
(3, 2, 12, '2026-01-19 11:30:00'),
(3, 4, 6, '2026-01-20 13:20:00'),
(4, 1, 20, '2026-01-21 08:45:00'),
(4, 3, 15, '2026-01-22 17:30:00'),
(5, 2, 7, '2026-01-23 10:15:00'),
(5, 4, 9, '2026-01-24 14:50:00'),
(6, 1, 4, '2026-01-25 12:10:00'),
(6, 5, 18, '2026-01-26 15:25:00'),
(7, 3, 11, '2026-01-27 09:40:00'),
(7, 6, 14, '2026-01-28 18:05:00'),
(8, 2, 2, '2026-01-29 10:55:00'),
(8, 4, 13, '2026-01-30 13:35:00'),
(9, 1, 16, '2026-01-31 11:20:00'),
(9, 5, 10, '2026-02-01 15:40:00'),
(10, 3, 19, '2026-02-02 14:15:00'),
(10, 6, 17, '2026-02-03 16:25:00');
````

## 准备 YAML

[mysql.yaml](./mysql.yaml)

```` yaml
name: mysql
input:
  - minLevel: int
  - configIds: string
  - page: int
workflows:
  - workflow:
    - step: MysqlExecute
      in:
        - (( env["mysqlUri"] ))
        - select id, player_id, config_id, level, created_time from hero where level > (( input["minLevel"] )) and config_id in (( "("+input["configIds"]+")" )) order by id desc limit (( (ParseInt(input["page"])-1)*5 )), 5
    - step: MatrixApply
      name: tableData
      in:
        - prev.out
        - col
        - 4
        - TimeConvert(value, "" , "2006-01-02 15:04")
  - workflow:
    - step: MatrixNew
      in:
        - ","
        - 2, Superman
        - 4, Spider-Man
        - 5, Batman
    - step: MatrixVLookUp
      in:
        - tableData.out
        - 2
        - prev.out
        - 0
        - 1
    - step: MatrixInsert
      name: result
      in:
        - prev.out
        - row
        - 0
        - Id
        - PlayerId
        - ConfigId
        - Level
        - UnlockAt
output: result.out
````

## 注册 MysqlExecute

[step_mysql.go](./step_mysql.go)

```` go
import (
	"fmt"
	"log/slog"

	"github.com/leeyaf/yflow"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	yflow.RegistStep("MysqlExecute", sMysqlExecute)
}

// 执行 MySQL 查询
//
// 可以保证列的顺序与查询顺序一致，默认丢弃列名
//
// 基于 gorm.io/gorm 实现
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

## 编写 main 函数

[main.go](./main.go)

```` go
package main

import (
    "fmt"
    "os"

    "github.com/leeyaf/yflow"
)

func main() {
    yamlData, err := os.ReadFile("myflow.yaml")
    if err != nil {
        panic(err)
    }

    // 注册全局环境变量
    yflow.RegistEnv("mysqlUri", "root:123456@tcp(127.0.0.1:3306)/yflow?charset=utf8mb4&parseTime=True&timeout=30s&loc=Local")

    // yaml 中定义的 input, 可以用表达式获取值
    input := []string{
        "1",              // minLevel
        "1001,1002,1005", // configIds
        "2",              // page
    }

    // 创建工作流
    job, err := yflow.NewJob(string(yamlData), input)
    if err != nil {
      panic(err)
    }

    // 执行工作流
    matrixResult, err := job.Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println(matrixResult)
	fmt.Println(job.ExecutionLog())
	// fmt.Println(job.MemoryModelLog())
}
````

## 运行

```` bash
go run .
````

## 输出

```` log
INFO MysqlExecute sql="select id, player_id, config_id, level, created_time from hero where level > 1 and config_id in (2,3,4) order by id desc limit 5, 5"
Matrix(6, 5)
Id, PlayerId, ConfigId, Level, UnlockAt
9, 5, Superman, 7, 2026-01-23 10:15
8, 4, 3, 15, 2026-01-22 17:30
6, 3, Spider-Man, 6, 2026-01-20 13:20
5, 3, Superman, 12, 2026-01-19 11:30
4, 2, 3, 8, 2026-01-18 16:45
Job execution: mysql
Workflow0: MysqlExecute [PASS] -> MatrixApply(tableData) [PASS]
Workflow1: MatrixNew [PASS] -> MatrixVLookUp [PASS] -> MatrixInsert(result) [PASS]
````