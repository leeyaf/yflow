package yflow_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/leeyaf/yflow"
)

const functionTestYaml = `
input:
- hello: string
workflows:
  - workflow:
    - step: MatrixNew
      in:
      - a b c d e f
      - g h i j k l
      - m n o p q r
      - 1 2 3 4 5 6
      - -1 -2 -3 -4 -5 -6
    - step: Test
      name: test
      in:
      - ${ Cell("prev.out", 0, 0) }
      - ${Row("prev.out",1)}
      - ${Col("prev.out",2)}
      - ${ParseInt(env.aInt)*2+100/2}
      - ${ParseFloat(env.aFloat)/2}
      - ${ParseFloat(env.aInt)*-2}
      - ${Join(Col("prev.out",0),",","'")}
      - ${Join(Row("prev.out",0),";","")}
      - ${Split(Cell("prev.in",4,0),",")}
      - ${TimeConvert(env.aTime,"","2006-01-02")}
      - ${Sprintf(input.hello,"Tom")}
      - ${Arange(0,10,2)}
      - ${Arange(5,-5,-3)}
output: test.out
`

func TestFunctions(t *testing.T) {
	yflow.RegistEnv("aInt", "-2")
	yflow.RegistEnv("aFloat", "-2.3")
	yflow.RegistEnv("aTime", fmt.Sprintf("%v", time.Now().Local()))
	job, err := yflow.NewJob(functionTestYaml, []string{
		"hello %v",
	})
	if err != nil {
		panic(err)
	}
	result, err := job.Execute()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
