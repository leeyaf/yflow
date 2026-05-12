package yflow

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"gopkg.in/yaml.v3"
)

const (
	MaxWaitStepSeconds = 30
)

var (
	steps     = make(map[string]*StepRegister)
	envMap    = make(map[string]string)
	functions = make(map[string]func(step *Step) any)
)

type StepRegister struct {
	Name   string
	Params []string
	f      func(step *Step, in *Matrix, out *Matrix)
}

func RegistFunction(key string, f any) {
	RegistFunctionWithStep(key, func(s *Step) any {
		return f
	})
}

func RegistFunctionWithStep(key string, f func(step *Step) any) {
	if key == "" {
		panic("empty key")
	}
	if f == nil {
		panic("no func")
	}
	if _, exists := functions[key]; exists {
		panic("repeat key " + key)
	}
	functions[key] = f
}

func RegistEnv(key string, val string) {
	if key == "" {
		panic("empty key")
	}
	if val == "" {
		panic("empty val")
	}
	if _, exists := envMap[key]; exists {
		panic("repeat key " + key)
	}
	envMap[key] = val
}

func GetEnv(key string) string {
	return envMap[key]
}

func RegistStep(name string, params []string, f func(step *Step, in *Matrix, out *Matrix)) {
	if name == "" {
		panic("no step name")
	}
	if params == nil {
		panic("no params")
	}
	if f == nil {
		panic("f is nil")
	}
	if _, exists := steps[name]; exists {
		panic("repeatname " + name)
	}
	steps[name] = &StepRegister{
		Name:   name,
		Params: params,
		f:      f,
	}
}

type Step struct {
	job        *Job
	workflow   *Workflow
	definition *StepDefinition
	register   *StepRegister
	in         *Matrix
	out        *Matrix
	done       chan (struct{})
}

func newStep(job *Job, workflow *Workflow, difinition *StepDefinition, register *StepRegister) *Step {
	return &Step{
		job:        job,
		workflow:   workflow,
		definition: difinition,
		register:   register,
		done:       make(chan struct{}),
	}
}

func (s *Step) Apply(value string, expression string) string {
	s.job.gengineCtx.Add("value", value)
	return s.job.Gengine(s, expression)
}

func (s *Step) GetDefinition() *StepDefinition {
	return s.definition
}

func (s *Step) GetRegister() *StepRegister {
	return s.register
}

func (s *Step) GetMatrix(exp string) *Matrix {
	return s.job.GetMatrix(s.definition, exp)
}

func (s *Step) getOut() *Matrix {
	select {
	case <-s.done: // chan 被关闭后，可以立即读到数据
		return s.out
	case <-time.After(MaxWaitStepSeconds * time.Second):
		panic("overtime")
	}
}

type Workflow struct {
	job   *Job
	steps []*Step
}

func (w *Workflow) execute() {
	for _, step := range w.steps {
		noExps := make([]string, 0)
		for _, row := range step.definition.In {
			data := row
			re := regexp.MustCompile(`\(\((.*?)\)\)`)
			matches := re.FindAllStringSubmatch(row, -1)
			for _, match := range matches {
				result := w.job.Gengine(step, match[1])
				data = strings.Replace(data, match[0], result, 1)
			}
			noExps = append(noExps, data)
		}

		in := NewMatrix()
		in.SetCol(0, noExps)
		step.in = in

		out := NewMatrix()
		step.out = out

		step.register.f(step, in, out)

		close(step.done)
	}
}

type Job struct {
	definition *JobDefinition
	input      map[string]string
	env        map[string]string
	workflows  []*Workflow
	gengineCtx *context.DataContext
}

/*
input:
(n, 2)
age, 18
gender, 1
score, 80
*/
func NewJob(yamlString string, input []string) *Job {
	definition := NewJobDefinition([]byte(yamlString))
	job := &Job{
		definition: definition,
		workflows:  make([]*Workflow, 0, len(definition.Workflows)),
		env:        envMap,
		gengineCtx: context.NewDataContext(),
	}

	inputMap := make(map[string]string)
	for i, row := range definition.Input {
		for k := range row {
			inputMap[k] = input[i]
			break
		}
	}
	job.input = inputMap

	job.gengineCtx.Add("input", inputMap)
	job.gengineCtx.Add("env", envMap)

	for _, workflowDef := range definition.Workflows {
		workflow := &Workflow{
			job:   job,
			steps: make([]*Step, 0, len(workflowDef.Steps)),
		}
		for _, stepDef := range workflowDef.Steps {
			stepReg := steps[stepDef.Step]
			if stepReg == nil {
				panic("no step " + stepDef.Step)
			}

			step := newStep(job, workflow, stepDef, stepReg)
			workflow.steps = append(workflow.steps, step)
		}
		job.workflows = append(job.workflows, workflow)
	}

	return job
}

func (j *Job) Execute() *Matrix {
	wg := new(sync.WaitGroup)
	wg.Add(len(j.workflows))
	for _, w := range j.workflows {
		go func(workflow *Workflow) {
			workflow.execute()
			wg.Done()
		}(w)
	}
	wg.Wait()

	return j.GetMatrix(nil, j.definition.Output)
}

func (j *Job) GetMatrix(stepDef *StepDefinition, exp string) *Matrix {
	seps := strings.Split(exp, ".")
	if len(seps) != 2 {
		panic("wrong name " + exp)
	}
	stepName := seps[0]
	matrixName := seps[1]

	var foundStep *Step
	if stepName == "prev" {
		if stepDef == nil {
			panic("no step but require prev")
		}
		for _, workflow := range j.workflows {
			for i, step := range workflow.steps {
				if step.definition == stepDef {
					foundStep = workflow.steps[i-1]
					break
				}
			}
			if foundStep != nil {
				break
			}
		}
	} else {
		for _, workflow := range j.workflows {
			for _, step := range workflow.steps {
				if step.definition.Name == stepName {
					foundStep = step
					break
				}
			}
			if foundStep != nil {
				break
			}
		}
	}

	if foundStep == nil {
		panic("not found step " + stepName)
	}

	if matrixName == "in" {
		return foundStep.in
	} else if matrixName == "out" {
		return foundStep.getOut()
	} else {
		panic("wrong matrix name " + matrixName)
	}
}

func (j *Job) Gengine(step *Step, expression string) string {
	for name, f := range functions {
		j.gengineCtx.Add(name, f(step))
	}

	script := "rule \"job\"\n" +
		"begin\n" +
		"return " + expression + "\n" +
		"end"
	rb := builder.NewRuleBuilder(j.gengineCtx)
	if err := rb.BuildRuleFromString(script); err != nil {
		panic(err)
	}
	eng := engine.NewGengine()
	if err := eng.Execute(rb, false); err != nil {
		panic(err)
	}
	result, err := eng.GetRulesResultMap()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v", result["job"])
}

type StepDefinition struct {
	Step string   `yaml:"step"`
	Name string   `yaml:"name"`
	In   []string `yaml:"in"`
}

type JobDefinition struct {
	Name      string              `yaml:"name"`
	Input     []map[string]string `yaml:"input"`
	Workflows []struct {
		Steps []*StepDefinition `yaml:"workflow"`
	} `yaml:"workflows"`
	Output string `yaml:"output"`
}

func NewJobDefinition(raw []byte) *JobDefinition {
	result := &JobDefinition{}
	if err := yaml.Unmarshal(raw, result); err != nil {
		panic(err)
	}
	return result
}
