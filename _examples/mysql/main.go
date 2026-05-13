package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/leeyaf/yflow"
)

func main() {
	yamlData, err := os.ReadFile("mysql.yaml")
	if err != nil {
		panic(err)
	}

	// 注册全局环境变量
	yflow.RegistEnv("mysqlUri", "root:123456@tcp(127.0.0.1:3306)/yflow?charset=utf8mb4&parseTime=True&timeout=30s&loc=Local")

	// yaml 中定义的 input, 可以用表达式获取值
	input := []string{
		"1",     // minLevel
		"2,3,4", // configIds
		"2",     // page
	}

	// 创建工作流
	job, err := yflow.NewJob(string(yamlData), input)
	if err != nil {
		panic(err)
	}

	// 执行工作流
	matrixResult, err := job.Execute()
	if err != nil {
		slog.Error(err.Error())
	}

	fmt.Println(matrixResult)
	fmt.Println(job.ExecutionLog())
	// fmt.Println(job.MemoryModelLog())
}
