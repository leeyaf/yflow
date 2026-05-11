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
	yflow.RegistEnv("mysqlUri", "root:123456@tcp(127.0.0.1:3306)/yamlql?charset=utf8mb4&parseTime=True&timeout=30s&loc=Local")

	// yaml 中定义的 input, 可以用表达式获取值
	input := []string{
		"1",                        // minLevel
		"1001,1002,1003,1004,1005", // configIds
		"2",                        // page
	}

	// 创建工作流
	job := yflow.NewJob(string(yamlData), input)

	// 执行工作流
	matrixResult := job.Execute()
	fmt.Println(matrixResult)
}
