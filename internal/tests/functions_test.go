package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/leeyaf/yflow"
)

func TestFunctions(t *testing.T) {
	raw, err := os.ReadFile("functions_test.yaml")
	if err != nil {
		t.Error(err)
	}

	yflow.RegistEnv("aInt", "-2")
	yflow.RegistEnv("aFloat", "-2.3")
	job, err := yflow.NewJob(string(raw), []string{})
	if err != nil {
		panic(err)
	}
	result, err := job.Execute()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
func TestTimePases(t *testing.T) {
	res, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", "2025-06-04 10:46:47 +0800 CST")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res.Format("2006-01-02 15:04:05"))
}
