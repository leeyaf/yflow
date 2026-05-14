# YFlow (Yaml Workflow Engine)

[简体中文](./README-zh.md)

A declarative workflow engine based on YAML configuration, allowing you to define complex data processing workflows through YAML, helping you turn Dirty Work into Clean Work.

## Features

* 🚀 **Declarative Configuration**: Define data processing workflows using YAML.
* 🔌 **Multiple Data Sources**: Easily connect to various data sources with just one Step.
* 📊 **Data Transformation**: Provides matrix operation functions (Transpose, VLookUp, Apply, etc.).
* ⚡ **Dynamic Expressions**: Based on [bilibili/gengine](https://github.com/bilibili/gengine), can be used within ``step.in``.
* 🔄 **Multi-Workflow Support**: Executes multiple workflows in parallel by default.
* 🧩 **Excellent Extensibility**: Both Steps and Functions can be extended, making it easy to customize your own processing logic.

## Quick Start

### Installation

```` bash
go get github.com/leeyaf/yflow
````

### Example

[A Relatively Complex Example](./_examples/mysql/README.md)

## Configuration Guide

### Basic Structure

```` yaml
name: myflow
input: # Original input
  - age: int # Type definition is not related to the engine; it can be used by the frontend (web interface)
  - phone: phoneNumber
workflows:
  - workflow: # First workflow
    - step: Test # Name of a registered Step
      name: result # Step name (optional)
      in: # Input for this Step
        - a
        - b
  - workflow: # Second workflow (executed in parallel)
    - step: Test
output: result.out # Points to a Step's output
````

The Steps registered by default are in [registed_step.go](./registed_step.go).

More Step implementations:

* [MySQL](./_examples/mysql/main.go)

### Expression Syntax

Expressions can be used in all ``step.in`` fields.

* Access input: ``(( input["paramName"] ))``
* Access env: ``(( env["paramName"] ))``
* Access the previous Step: ``(( Cell("prev.out", 0, 0) ))``
* Nested usage: ``(( ParseInt(input["paramName"]) ))`` (The **final** return result of the expression will be formatted as a ``string``).
* Simple operations: ``(( ParseInt(input["paramName"])*10 ))``

Refer to [bilibili/gengine](https://github.com/bilibili/gengine/wiki/语法) for expression details.

The Functions registered by default are in [registed_function.go](./registed_function.go).

## Development Guide

### Adding a Step

```` go
func sMyStep(s *yflow.Step, in *yflow.Matrix, out *yflow.Matrix) {
    // Write processing results to 'out'
}

yflow.RegistStep("MyStep", sMyStep)
````

### Adding a Function

```` go
// For functions that do NOT need access to the current Step
func fMyFunction(value string) string {
    return "hello " + value
}

yflow.RegistFunction("MyFunction", fMyFunction)

// For functions that DO need access to the current Step
func fMyFunction2(step *yflow.Step) any {
    return func(matrixName string) string {
        rows, cols := step.GetMatrix(matrixName).Shape()
        return fmt.Sprintf("%v, %v", rows, cols)
    }
}

yflow.RegistFunctionWithStep("MyFunction2", fMyFunction2)
````

## Memory Model

![memory-model](./memory-model.png)

## Background

Previously, when I served as the lead server programmer at a gaming company, after the game launched, operations and planning teams continuously requested data and reports from the server. Among these requests, some were one-time, while others were recurring. My initial strategy was simple: evaluate the nature of the request (short-term vs. long-term). Long-term needs were implemented in code, and one-time requests were handled with SQL + Excel. However, the problem was: many supposedly one-time requests would suddenly reappear one day. Recalling the lengthy steps I had previously implemented, I wondered: Could there be a low-cost way to formalize these lengthy procedures?

Thus, I conceived and developed the original version (internal company code ``fakegm``). With ``fakegm``, we met various demands at a very low cost and no longer needed to evaluate their long-term or short-term nature (treating all as long-term). As it evolved, our server-side business code was eventually reduced to just login, timers, and other essentials; everything else became implementations of various Steps. The codebase shrank, but functionality became richer. We integrated numerous data sources: MySQL, Redis, Kubernetes, AliyunOSS, AliyunImage, GrafanaCloud, internal configuration management platforms, and internal game servers. On the frontend (web interface), aside from the basic Job management and execution page, we also created dedicated, polished pages for some particularly high-frequency needs. Of course, the basic Job management page could already fulfill all requirements, but the customized pages struck the perfect balance for various stakeholders.

After accomplishing all this, we even began to look forward to various requests whose nature (long-term or short-term) was unclear: when existing Steps could meet the need, writing a YAML file would deliver the result; when existing Steps fell short, it was an opportunity to add a new one to our Step library!

Now, ``yflow`` inherits the purest essence of ``fakegm``, completely recoded and redesigned to make Clean Work even cleaner.

## Contribution

Steps and Functions that are generally useful, fundamental, and not biased towards specific projects are registered by default. Those with a specific bias should be placed in ``/_examples/``.

Issues and Pull Requests are welcome!

1. Fork the project
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License