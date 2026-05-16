package main_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/leeyaf/yflow"
)

func TestCurd(t *testing.T) {
	yamlData, err := os.ReadFile("mysql_curd.yaml")
	if err != nil {
		panic(err)
	}

	yflow.RegistEnv("mysqlUri", "root:123456@tcp(127.0.0.1:3306)/yflow?charset=utf8mb4&parseTime=True&timeout=30s&loc=Local")

	input := []string{
		"6",
		"1",
		"101",
	}

	job, err := yflow.NewJob(string(yamlData), input)
	if err != nil {
		panic(err)
	}
	matrixResult, err := job.Execute()
	if err != nil {
		slog.Error(err.Error())
	}
	fmt.Println(matrixResult)
	fmt.Println(job.ExecutionLog())
	fmt.Println(job.MemoryModelLog())
}

func TestParallel(t *testing.T) {
	yamlData, err := os.ReadFile("mysql_parallel.yaml")
	if err != nil {
		panic(err)
	}

	yflow.RegistEnv("mysqlUriLayout", "root:123456@tcp(127.0.0.1:3306)/yflow%v?charset=utf8mb4&parseTime=True&timeout=30s&loc=Local")

	input := []string{
		"1", // minLevel
	}

	job, err := yflow.NewJob(string(yamlData), input)
	if err != nil {
		panic(err)
	}
	matrixResult, err := job.Execute()
	if err != nil {
		slog.Error(err.Error())
	}
	fmt.Println(matrixResult)
	fmt.Println(job.ExecutionLog())
	fmt.Println(job.MemoryModelLog())
}
