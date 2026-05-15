package redis

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-redis/redis/v8"
	"github.com/google/shlex"
	"github.com/leeyaf/yflow"
)

func init() {
	yflow.RegistStep("RedisExecute", sRedisExecute)
}

// 执行 Redis 命令
//
// 基于 go-redis 和 shlex
//
// 结果是数组时自动放在矩阵的第一列
func sRedisExecute(s *yflow.Step, in *yflow.Matrix, out *yflow.Matrix) {
	uri := in.Get(0, 0)     // Redis 链接地址
	command := in.Get(1, 0) // 命令

	parts, err := shlex.Split(command)
	if err != nil {
		panic(err)
	}
	doArgs := make([]any, 0, len(parts))
	for _, part := range parts {
		doArgs = append(doArgs, part)
	}

	options, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}
	cli := redis.NewClient(options)
	defer cli.Close()

	result, err := cli.Do(context.Background(), doArgs...).Result()
	if err != nil {
		panic(err)
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.Slice:
		for i, data := range result.([]any) {
			out.Set(i, 0, fmt.Sprintf("%v", data))
		}
	default:
		out.Set(0, 0, fmt.Sprintf("%v", result))
	}
}
