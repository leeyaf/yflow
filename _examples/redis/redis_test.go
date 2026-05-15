package redis_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/leeyaf/yflow"
)

func TestRedis(t *testing.T) {
	yamlData, err := os.ReadFile("redis_test.yaml")
	if err != nil {
		panic(err)
	}

	yflow.RegistEnv("redisUri", "redis://yflow:123455@localhost:6379/0?dial_timeout=3&read_timeout=3&write_timeout=3")

	input := []string{
		"Tom",
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
