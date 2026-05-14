package yflow

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/bilibili/gengine/builder"
	gcontext "github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"gopkg.in/yaml.v3"
)

var envMap = make(map[string]string)

func RegistEnv(key string, value string) error {
	if key == "" {
		return errors.New("empty key")
	}
	if value == "" {
		return errors.New("empty value")
	}
	if _, exists := envMap[key]; exists {
		slog.Warn("duplicate", "key", key)
	}

	envMap[key] = value
	return nil
}

func GetEnv(key string) string {
	return envMap[key]
}

var functions = make(map[string]func(step *Step) any)

func RegistFunction(key string, f any) {
	RegistFunctionWithStep(key, func(s *Step) any {
		return f
	})
}

func RegistFunctionWithStep(key string, f func(step *Step) any) error {
	if key == "" {
		return errors.New("empty key")
	}
	if f == nil {
		return errors.New("f is nil")
	}
	if _, exists := functions[key]; exists {
		slog.Warn("duplicate", "key", key)
	}

	functions[key] = f
	return nil
}

var steps = make(map[string]*StepRegister)

type StepRegister struct {
	Name string
	f    func(step *Step, in *Matrix, out *Matrix)
}

func RegistStep(name string, f func(step *Step, in *Matrix, out *Matrix)) error {
	if name == "" {
		return errors.New("empty name")
	}
	if f == nil {
		return errors.New("f is nil")
	}
	if _, exists := steps[name]; exists {
		slog.Warn("duplicate", "name", name)
	}

	steps[name] = &StepRegister{
		Name: name,
		f:    f,
	}
	return nil
}

type StepError struct {
	step    *Step
	message any
	stack   []byte
}

func (e *StepError) Error() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf("%v [%v]", e.message, e.step.pathLog()))
	b.WriteString(fmt.Sprintf("\n%v", string(e.stack)))
	return b.String()
}

const StepMaxLifetime = 30 * time.Second

type stepState uint8

const (
	stateDefault = stepState(0) + iota
	statePass
	stateFail
	stateCancelOrTimeout
)

type Step struct {
	workflow   *Workflow
	definition *StepDefinition
	register   *StepRegister
	in         *Matrix
	inDone     chan (struct{})
	out        *Matrix
	outDone    chan (struct{})
	state      stepState
}

func newStep(workflow *Workflow, difinition *StepDefinition, register *StepRegister) *Step {
	step := &Step{
		workflow:   workflow,
		definition: difinition,
		register:   register,
		inDone:     make(chan struct{}),
		outDone:    make(chan struct{}),
	}
	return step
}

func (s *Step) pathLog() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf("Workflow%v", s.workflow.index))
	for _, step := range s.workflow.steps {
		name := step.definition.Name
		if len(name) != 0 {
			name = "(" + name + ")"
		}

		b.WriteString(fmt.Sprintf(" -> %v%v", step.definition.Step, name))
		if step == s {
			break
		}
	}
	return b.String()
}

// 获取其他 Step 的 IO 矩阵，
// 如果 Workflow 还未执行到其他 Step，会等待
//
// yourStepName.in yourStepName.out
//
// prev.in prev.out
func (s *Step) GetMatrix(expression string) *Matrix {
	stepName, matrixName, err := s.workflow.job.parseMatrixNameExp(expression)
	if err != nil {
		panic(err)
	}

	if stepName == s.definition.Name {
		panic("cant use this function to access itself")
	}

	if stepName == "prev" {
		prevStep, err := s.workflow.getPrevStep()
		if err != nil {
			panic(err)
		}
		if matrixName == "in" {
			return prevStep.in
		} else {
			return prevStep.out
		}
	} else {
		selfIndex := func() int {
			for i, step := range s.workflow.steps {
				if step == s {
					return i
				}
			}
			panic("memory exception")
		}()

		// 先从所属的 workflow 查找
		for i := 0; i < selfIndex && i < len(s.workflow.steps); i++ {
			step := s.workflow.steps[i]
			if step.definition.Name == stepName {
				if matrixName == "in" {
					return step.in
				} else {
					return step.out
				}
			}
		}

		// 再从其他 workflow 查找
		return func() *Matrix {
			for _, workflow := range s.workflow.job.workflows {
				if workflow == s.workflow {
					continue
				}
				for _, step := range workflow.steps {
					if step.definition.Name == stepName {
						if matrixName == "in" {
							<-step.inDone
							return step.in
						} else {
							<-step.outDone
							return step.out
						}
					}
				}
			}
			panic(fmt.Errorf("not found matrix: %v", expression))
		}()
	}
}

func (s *Step) execute(ctx context.Context) {
	noExps := make([]string, 0)
	for _, row := range s.definition.In {
		select {
		case <-ctx.Done():
			return
		default:
		}

		noExp := row
		re := regexp.MustCompile(`\(\((.*?)\)\)`)
		matches := re.FindAllStringSubmatch(row, -1)
		for _, match := range matches {
			result := s.GengineExecute(match[1])
			noExp = strings.Replace(noExp, match[0], result, 1)
		}
		noExps = append(noExps, noExp)
	}

	in := NewMatrix()
	in.SetCol(0, noExps)
	s.in = in
	close(s.inDone)

	out := NewMatrix()
	s.out = out

	select {
	case <-ctx.Done():
		return
	default:
		s.register.f(s, in, out)
		close(s.outDone)
	}
}

func (s *Step) GengineAddContext(key string, obj any) {
	s.workflow.gengineCtx.Add(key, obj)
}

func (s *Step) GengineExecute(expression string) string {
	// 更新 Function 中的 Step 指针
	for name, f := range functions {
		s.workflow.gengineCtx.Add(name, f(s))
	}

	result, err := s.workflow.gengineExecute(expression)
	if err != nil {
		panic(err)
	}
	return result
}

type Workflow struct {
	job        *Job
	steps      []*Step
	gengineCtx *gcontext.DataContext
	stepIndex  int
	index      int
}

func newWorkflow(job *Job, definition *WorkflowDefinition, workflowIndex int) (*Workflow, error) {
	workflow := &Workflow{
		job:        job,
		steps:      make([]*Step, 0, len(definition.Steps)),
		gengineCtx: gcontext.NewDataContext(),
		index:      workflowIndex,
	}

	workflow.gengineCtx.Add("input", job.input)
	workflow.gengineCtx.Add("env", envMap)

	for _, stepDef := range definition.Steps {
		registedStep := steps[stepDef.Step]
		if registedStep == nil {
			return nil, fmt.Errorf("definition not found: %v", stepDef.Step)
		}

		step := newStep(workflow, stepDef, registedStep)
		workflow.steps = append(workflow.steps, step)
	}

	return workflow, nil
}

func (w *Workflow) executeStep(ctx context.Context, s *Step) error {
	stepCtx, cancel := context.WithTimeout(ctx, StepMaxLifetime)
	defer cancel()

	errorChan := make(chan error, 1)
	go func(step *Step) {
		defer func() {
			if r := recover(); r != nil {
				err := &StepError{
					step:    step,
					message: r,
					stack:   debug.Stack(),
				}
				errorChan <- err
			} else {
				errorChan <- nil
			}
		}()
		step.execute(stepCtx)
	}(s)

	select {
	case err := <-errorChan:
		s.state = statePass
		if err != nil {
			s.state = stateFail
		}
		return err
	case <-stepCtx.Done():
		s.state = stateCancelOrTimeout
		return &StepError{
			step:    s,
			message: "cancel or timeout",
			stack:   debug.Stack(),
		}
	}
}

// 顺序执行 Step
func (w *Workflow) execute(ctx context.Context) error {
	for i, s := range w.steps {
		select {
		case <-ctx.Done():
			return nil
		default:
			w.stepIndex = i

			if err := w.executeStep(ctx, s); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Workflow) getPrevStep() (*Step, error) {
	if w.stepIndex == 0 {
		return nil, errors.New("cant call prev here")
	}
	return w.steps[w.stepIndex-1], nil
}

func (w *Workflow) gengineExecute(expression string) (string, error) {
	script := "rule \"job\"\n" +
		"begin\n" +
		"return " + expression + "\n" +
		"end"
	rb := builder.NewRuleBuilder(w.gengineCtx)
	if err := rb.BuildRuleFromString(script); err != nil {
		return "", err
	}
	eng := engine.NewGengine()
	if err := eng.Execute(rb, false); err != nil {
		return "", err
	}
	result, err := eng.GetRulesResultMap()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result["job"]), nil
}

type Job struct {
	definition *JobDefinition
	input      map[string]string
	env        map[string]string
	workflows  []*Workflow
}

func NewJob(yamlString string, input []string) (*Job, error) {
	definition, err := NewJobDefinition([]byte(yamlString))
	if err != nil {
		return nil, err
	}

	job := &Job{
		definition: definition,
		workflows:  make([]*Workflow, 0, len(definition.Workflows)),
		env:        envMap,
	}

	inputMap := make(map[string]string)
	for i, row := range definition.Input {
		for k := range row {
			inputMap[k] = input[i]
			break
		}
	}
	job.input = inputMap

	for index, workflowDef := range definition.Workflows {
		workflow, err := newWorkflow(job, workflowDef, index)
		if err != nil {
			return nil, err
		}
		job.workflows = append(job.workflows, workflow)
	}

	return job, nil
}

func (j *Job) ExecutionLog() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf("Job execution: %v", j.definition.Name))
	for i, workflow := range j.workflows {
		b.WriteString(fmt.Sprintf("\nWorkflow%v: ", i))
		for j, step := range workflow.steps {
			name := step.definition.Name
			if len(name) != 0 {
				name = "(" + name + ")"
			}

			state := ""
			if step.state == statePass {
				state = " [PASS]"
			} else if step.state == stateFail {
				state = " [FAIL]"
			} else if step.state == stateCancelOrTimeout {
				state = " [CANCEL or TIMEOUT]"
			}

			if j > 0 {
				b.WriteString(" -> ")
			}
			b.WriteString(step.definition.Step + name)
			if state != "" {
				b.WriteString(state)
			}
		}
	}
	return b.String()
}

func (j *Job) MemoryModelLog() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf("Job memory model: %v", j.definition.Name))
	b.WriteString(fmt.Sprintf("\ninput: %v", j.input))
	b.WriteString(fmt.Sprintf("\noutput: %v", j.definition.Output))
	for _, workflow := range j.workflows {
		for _, step := range workflow.steps {
			if step.in != nil {
				b.WriteString(fmt.Sprintf("\n%v => in:\n%v", step.pathLog(), step.in))
			}
			if step.out != nil {
				b.WriteString(fmt.Sprintf("\n%v => out:\n%v", step.pathLog(), step.out))
			}
		}
	}
	return b.String()
}

// Job 的最小执行单位是 Step，
// Workflow 会捕获在 Step 中产生的 panic，
// 所以在 Step、Function 的实现里可以直接抛出 panic。
//
// 这样做的目的是降低 Step、Function 的开发成本，
// 简化 YAML 的定义。
//
// 工作机制：
//
// 1. 并行执行多 Workflow
//
// 2. Workflow 在捕获到 panic 后立即通知其他 Workflow 停止执行
//
// 3. 创建包含错误描述、Step 路径、goroutine 堆栈的 StepError 对象，
// 传递给 Job 返回
//
// 4. Step 执行带有超时机制，超时时抛出 StepError，
// StepMaxLifetime = 30 seconds
func (j *Job) Execute() (out *Matrix, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(j.workflows))

	errorChan := make(chan error, 1)
	for _, w := range j.workflows {
		go func(workflow *Workflow) {
			defer wg.Done()

			if err := workflow.execute(ctx); err != nil {
				select {
				case errorChan <- err:
				default:
					// 只捕获一个 panic
				}
			}
		}(w)
	}

	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case err := <-errorChan:
		// 通知其他 workflow 停止
		cancel()
		<-doneChan

		return nil, err
	case <-doneChan:
		return j.GetMatrix(j.definition.Output)
	}
}

func (j *Job) parseMatrixNameExp(expression string) (stepName, matrixName string, err error) {
	parts := strings.Split(expression, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("illegal expression: %v", expression)
	}

	stepName = parts[0]
	matrixName = parts[1]
	if matrixName != "in" && matrixName != "out" {
		return "", "", fmt.Errorf("illegal matrix name: %v", matrixName)
	}
	return stepName, matrixName, nil
}

func (j *Job) GetMatrix(expression string) (*Matrix, error) {
	stepName, matrixName, err := j.parseMatrixNameExp(expression)
	if err != nil {
		return nil, err
	}

	return func() (*Matrix, error) {
		for _, workflow := range j.workflows {
			for _, step := range workflow.steps {
				if step.definition.Name == stepName {
					if matrixName == "in" {
						return step.in, nil
					} else {
						return step.out, nil
					}
				}
			}
		}
		return nil, fmt.Errorf("not found matrix: %v", expression)
	}()
}

type StepDefinition struct {
	Step string   `yaml:"step"`
	Name string   `yaml:"name"`
	In   []string `yaml:"in"`
}

type WorkflowDefinition struct {
	Steps []*StepDefinition `yaml:"workflow"`
}

type JobDefinition struct {
	Name      string                `yaml:"name"`
	Input     []map[string]string   `yaml:"input"`
	Workflows []*WorkflowDefinition `yaml:"workflows"`
	Output    string                `yaml:"output"`
}

func NewJobDefinition(raw []byte) (*JobDefinition, error) {
	result := &JobDefinition{}
	if err := yaml.Unmarshal(raw, result); err != nil {
		return nil, err
	}
	return result, nil
}
