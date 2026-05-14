package yflow_test

import (
	"fmt"
	"testing"

	"github.com/leeyaf/yflow"
)

const stepTestYaml = `
workflows:
  - workflow:
    - step: Test
      name: test
      in:
      - hello (( env["engineName"] ))
    - step: MatrixNew
      name: bestOne
      in:
      - ","
      - a,b,c,d,e,f
      - g,h,i,j,k,l
      - m,n,o,p,q,r
      - 10,2,73,94,5,670
      - -1,-2,-3,-4,-5,-6
    - step: MatrixNew
      name: renameBestOne
      in:
      - ","
      - a, A
      - b, B
      - c, C
      - d, D
    - step: MatrixEmpty
      in:
      - 3
      - 4
      - SuperCool
    - step: MatrixTranspose
      name: transposeBestOne
      in:
      - bestOne.out
    - step: MatrixSort
      in:
      - prev.out
      - 3
      - asc
    - step: MatrixInsert
      in:
      - bestOne.out
      - col
      - 1
      - i0
      - i1
      - i2
      - i3
      - i4
      - i5
    - step: MatrixSelect
      in:
      - bestOne.out
      - row
      - 0
      - 2
    - step: MatrixVLookUp
      in:
      - transposeBestOne.out
      - 0
      - renameBestOne.out
      - 0
      - 1
    - step: MatrixApply
      in:
      - bestOne.out
      - row
      - 4
      - value
      - ParseInt(value)*-1
    - step: MatrixMerge
      in:
      - bestOne.out
      - renameBestOne.out
      - col
output: test.out
`

func TestSteps(t *testing.T) {
	yflow.RegistEnv("engineName", "leeyaf/yflow")
	job, err := yflow.NewJob(stepTestYaml, []string{})
	if err != nil {
		panic(err)
	}
	result, err := job.Execute()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	fmt.Println(job.MemoryModelLog())
}
