# YFlow (Yaml Workflow Engine)

一个基于 YAML 配置的声明式工作流引擎，允许你通过 YAML 定义复杂的数据处理流程，帮你把 Ditry Work 变成 Clean Work。

## 特性

* 🚀 **声明式配置**：使用 YAML 定义数据处理工作流
* 🔌 **多数据源**：只要一个 Step，就可以轻松接入各种数据源
* 📊 **数据转换**：提供矩阵操作函数（Transpose、VLookUp、Apply 等）
* ⚡ **动态表达式**：基于 [bilibili/gengine](https://github.com/bilibili/gengine) 实现，可以在 ``step.in`` 中使用
* 🔄 **多工作流支持**：默认并行执行多个工作流
* 🧩 **良好扩展性**：Step、Function 都可以扩展，轻松定制自己的处理逻辑

## 背景

之前我在一家游戏公司做服务器主程时，游戏上线后，运营、策划陆续开始向服务器索要数据、报表。在这些需求里，有部分是一次性的，有部分是长期的，我当时的策略很简单：评估需求的长短期属性，长期的 Code 落地，一次性的 SQL + Excel 处理。但问题是：有很多一次性的需求，突然在某一天又冒了出来，回想曾经实现的冗长步骤，我就在想：能不能找一种低成本的方式把冗长步骤落地呢？

于是，我构思并开发出了原始版本（公司内部代码 ``fakegm``）。借助 ``fakegm``，我们用很低的成本满足了各种需求，也不需要再评估长短期属性（全部视为长期）。发展到后期，我们服务端的业务代码就只剩下了登陆、定时器和其他必要的，其余的全是各种 Step 的实现。代码减少了，功能却更丰富了，我们打通了非常多的数据源：MySQL、Redis、Kubernetes、AliyunOSS、AliyunImage、GrafanaCloud、内部配置表管理平台、内部游戏服。而前端（网页端），除了基础的 Job 管理执行页面外，我们针对一些特别高频的需求，还定制了专用的漂亮页面，当然基础的 Job 管理执行页面已经可以实现所有的需求了，但定制的漂亮页面在各方的需求中找到了完美的平衡点。

当这一切做好后，我们甚至开始期待各种看不到长短期属性的需求：当已有的 Step 可以满足时，写一个 YAML 就可以交付了；当已有的 Step 无法满足时，又可以给我们的 Step 库加个新的了！

现在，``yflow`` 继承了 ``fakegm`` 最纯粹的思想，重新编码、重新设计，让 Clen Work 更 Clean。

## 快速开始

### 安装

```` bash
go get github.com/leeyaf/yflow
````

### 一个比较复杂的示例

#### 创建 MySQL 表 ``yflow.hero``

```` sql
+-----+-------+-----------+-----------+---------------------+
| id  | level | player_id | config_id | created_time        |
+-----+-------+-----------+-----------+---------------------+
| 470 |    62 |        19 |      1002 | 2025-08-10 10:48:12 |
| 469 |     3 |        15 |      1002 | 2025-05-26 10:48:12 |
| 468 |    88 |         2 |      1002 | 2025-09-21 10:48:12 |
| 467 |    52 |        11 |      1005 | 2026-04-11 10:48:12 |
| 466 |    30 |        13 |      1009 | 2025-08-29 10:48:12 |
+-----+-------+-----------+-----------+---------------------+
````

#### 创建 ``myflow.yaml``

```` yaml
name: myflow
input:
  - minLevel: int
  - configIds: string
  - page: int
workflows:
  - workflow:
    - step: MysqlExecute
      in:
        - (( env["mysqlUri"] ))
        - select id, level, player_id, config_id, created_time from hero where level > (( input["minLevel"] )) and config_id in (( "("+input["configIds"]+")" )) order by id desc limit (( (ParseInt(input["page"])-1)*5 )), 5
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
        - 1001, Superman
        - 1002, Spider-Man
        - 1003, Batman
    - step: MatrixVLookUp
      in:
        - tableData.out
        - 3
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
        - Level
        - PlayerId
        - ConfigId
        - UnlockAt
output: result.out
````

#### 运行

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

````

#### 输出

```` log
INFO MysqlExecute sql="select id, level, player_id, config_id, created_time from hero where level > 1 and config_id in (1001,1002,1003,1004,1005) order by id desc limit 5, 5"
(6, 5)
Id, Level, PlayerId, ConfigId, UnlockAt
464, 14, 18, 1004, 2026-01-15 10:48
460, 83, 16, Batman, 2026-04-01 10:48
459, 6, 7, 1004, 2025-05-28 10:48
457, 74, 3, Superman, 2025-09-29 10:48
453, 13, 18, Batman, 2025-08-06 10:48
````

## 配置说明

### 基本结构

```` yaml
name: myflow
input: # 原始输入
  - age: int # 类型定义与 engine 无关, 可以给前端（网页端）用
  - phone: phoneNumber
workflows:
  - workflow: # 第一个工作流
    - step: MatrixNew # 已注册的 Step 名
      name: result # 步骤名称（可选）
      in: # 这个 Step 的输入
          - "," # MatrixNew 要求第 0 行 0 列是分隔符定义
        - a, b, c, d # 矩阵的第一行
        - 1, 2, 3, 4 # 矩阵的第二行
  - workflow: # 第二个工作流（并行执行）
    - step: Test
output: result.out # 指向 Step
````

可用的 Step [在这里](./registed-step.md)

更多的 Step 实现：

* [MySQL](./more/step-mysql.md)

### 表达式语法

可以在所有 ``step.in`` 中使用表达式

* 访问 input: ``(( input["paramName"] ))``
* 访问 env: ``(( env["paramName"] ))``
* 访问上一个 Step: ``(( Cell("prev.out", 0, 0) ))``
* 嵌套使用: ``(( ParseInt(input["paramName"]) ))``（表达式**最终的**返回结果会被格式化为 ``string``）
* 简单运算: ``(( ParseInt(input["paramName"])*10 ))``

参考 [bilibili/gengine语法](https://github.com/bilibili/gengine/wiki/语法)

可用的 Function [参考这里](./registed-function.md)

## 开发指南

### 添加 Step

```` go
func sMyStep(s *yflow.Step,in *yflow.Matrix, out *yflow.Matrix) {
    // 把处理结果写入到 out
}

yflow.RegistStep(
    "MyStep",
    []string{"Param1", "Param2", ".."},
    sMyStep,
)
````

### 添加 Function

#### 实现

```` go
// 不需要访问当前 Step 的
func fMyFunction(value string) string {
    return "hello " + value
}

yflow.RegistFunction("MyFunction", fMyFunction)

// 需要的
func fMyFunction2(step *yflow.Step) any {
    return func(matrixName string) string {
        rows, cols := step.GetMatrix(matrixName).Shape()
        return fmt.Sprintf("%v, %v", rows, cols)
    }
}

yflow.RegistFunctionWithStep("MyFunction2", fMyFunction2)
````

## 贡献

通用性较好的、基础的、无具体项目偏向的 Step 和 Function 默认注册，有偏向的请放在 ``/more/`` 里，可以拷贝源代码到项目中使用。

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
1. 创建功能分支
1. 提交更改
1. 推送到分支
1. 创建 Pull Request

## 许可证

MIT License