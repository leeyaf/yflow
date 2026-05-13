package tests

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/leeyaf/yflow"
)

func TestWait(t *testing.T) {
	doneChan := make(chan (struct{}))
	time.AfterFunc(3*time.Second, func() {
		close(doneChan)
	})

	go func() {
		slog.Info("g0 start wait")
		<-doneChan
		slog.Info("g0 done")
	}()

	time.AfterFunc(1*time.Second, func() {
		slog.Info("g1 start wait")
		<-doneChan
		slog.Info("g1 done")
	})

	time.AfterFunc(2*time.Second, func() {
		slog.Info("g2 start wait")
		<-doneChan
		slog.Info("g2 done")
	})

	time.AfterFunc(3*time.Second, func() {
		slog.Info("g3 start wait")
		<-doneChan
		slog.Info("g3 done")
	})

	time.AfterFunc(4*time.Second, func() {
		slog.Info("g4 start wait")
		<-doneChan
		slog.Info("g4 done")
	})

	time.AfterFunc(5*time.Second, func() {
		slog.Info("g5 start wait")
		<-doneChan
		slog.Info("g5 done")
	})

	time.Sleep(8 * time.Second)
}

func TestYaml(t *testing.T) {
	raw, err := os.ReadFile("job_test.yaml")
	if err != nil {
		t.Error(err)
	}

	yflow.RegistEnv("mysqlUri", "root:123456@tcp(127.0.0.1:3306)/yamlql?charset=utf8mb4&parseTime=True&timeout=30s&loc=Local")
	job, err := yflow.NewJob(string(raw), []string{"1", "1001,1002,1003,1004,1005", "1"})
	if err != nil {
		panic(err)
	}
	fmt.Println(job.Execute())
}
