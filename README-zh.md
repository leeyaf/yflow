# YFlow (Yaml Workflow Engine)

一个基于 YAML 配置的声明式工作流引擎，允许你通过 YAML 定义复杂的数据处理流程，帮你把 Ditry Work 变成 Clean Work。

## 特性

* 🚀 **声明式配置**：使用 YAML 定义数据处理工作流
* 🔌 **多数据源**：只要一个 ``Step``，就可以轻松接入各种数据源
* 📊 **数据转换**：提供丰富的矩阵操作函数
* ⚡ **函数表达式**：内置 [bilibili/gengine](https://github.com/bilibili/gengine) 规则引擎
* 🔄 **多工作流支持**：并行执行多个工作流
* 🧩 **良好扩展性**：轻松扩展 ``Step`` 和 ``Function``

## 快速开始

### 安装

```` bash
go get github.com/leeyaf/yflow
````

### 示例

[一个比较复杂的示例](./_examples/mysql/README.md)

## 配置说明

### 基本结构

```` yaml
name: myflow
input: # 原始输入
  - age: int # 类型定义与 engine 无关, 可以给前端（网页端）用
  - phone: phoneNumber
workflows:
  - workflow: # 第一个工作流
    - step: Test # 已注册的 Step 名
      name: result # 步骤名称（可选）
      in: # 输入矩阵
        - a # 第一行
        - b # 第二行
  - workflow: # 第二个工作流（并行执行）
    - step: Test
output: result.out # 指向 Step 的矩阵
````

### 函数表达式

在 YAML 里所有 ``Step`` 的 ``in`` 都可以使用 ``${YourExpression}``

| 表达式 | 说明 |
| ----- | --- |
| ${prev.in} | 上个 ``Step`` 的输入矩阵地址 |
| ${prev.out} | 上个 ``Step`` 的输出矩阵地址 |
| ${yourStepName.in} | 指定 ``Step`` 的输入矩阵地址 |
| ${input.age} | 获取 ``input`` 中的值 |
| ${env.mysqlUri} | 获取 ``env`` 中的值 |
| ${ParseInt("6")} | 调用 ``Function`` |
| ${10*3} | 简单数学元算，参考 [bilibili/gengine语法](https://github.com/bilibili/gengine/wiki/语法) |


* 直接可用的 ``Step`` 在 [registed_step.go](./registed_step.go)，示例 [registed_step_test.go](./registed_step_test.go)
* 更多 ``Step`` 实现
  * [MySQL](./_examples/mysql/step_mysql.go)
  * [Redis](./_examples/redis/step_redis.go)
* 直接可用的 ``Function`` 在 [registed_function.go](./registed_function.go)，示例 [registed_function_test.go](./registed_function_test.go)

### 其他

* 因为 ``Matrix`` 中的元素是 ``string`` 型的，所以函数表达式的输出会被转为 ``string``
* 在编写 YAML 时，如果一个字符串实际代表多个值，应使用空格分隔 ``shell-style``

## 开发指南

### 添加 ``Step``

```` go
func sMyStep(s *yflow.Step, in *yflow.Matrix, out *yflow.Matrix) {
    // 把处理结果写入到 out
}

yflow.RegistStep("MyStep", sMyStep)
````

### 添加 ``Function``

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

## 内存模型

![内存模型](./memory-model.png)

## 背景

之前我在一家游戏公司做服务器主程时，游戏上线后，运营、策划陆续开始向服务器索要数据、报表。在这些需求里，有部分是一次性的，有部分是长期的，我当时的策略很简单：评估需求的长短期属性，长期的 Code 落地，一次性的 SQL + Excel 处理。但问题是：有很多一次性的需求，突然在某一天又冒了出来，回想曾经实现的冗长步骤，我就在想：能不能找一种低成本的方式把冗长步骤落地呢？

于是，我构思并开发出了原始版本（公司内部代码 ``fakegm``）。借助 ``fakegm``，我们用很低的成本满足了各种需求，也不需要再评估长短期属性（全部视为长期）。发展到后期，我们服务端的业务代码就只剩下了登陆、定时器和其他必要的，其余的全是各种 ``Step`` 的实现。代码减少了，功能却更丰富了，我们打通了非常多的数据源：MySQL、Redis、Kubernetes、AliyunOSS、AliyunImage、GrafanaCloud、内部配置表管理平台、内部游戏服。而前端（网页端），除了基础的 ``Job`` 管理执行页面外，我们针对一些特别高频的需求，还定制了专用的漂亮页面，当然基础的 ``Job`` 管理执行页面已经可以实现所有的需求了，但定制的漂亮页面在各方的需求中找到了完美的平衡点。

当这一切做好后，我们甚至开始期待各种看不到长短期属性的需求：当已有的 ``Step`` 可以满足时，写一个 YAML 就可以交付了；当已有的 ``Step`` 无法满足时，又可以给我们的 ``Step`` 库加个新的了！

现在，``yflow`` 继承了 ``fakegm`` 最纯粹的思想，重新编码、重新设计，让 Clen Work 更 Clean。

## 贡献

通用性较好的、基础的、无具体项目偏向的 ``Step`` 和 ``Function`` 默认注册，有偏向的请放在 ``/_examples/`` 里。

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
1. 创建功能分支
1. 提交更改
1. 推送到分支
1. 创建 Pull Request

## 许可证

MIT License