package transfer

import "testing"

type Src struct {
	Name string
	Age  int
	Info []string
}

type Dst struct {
	Name  string
	Age   int
	Info  []string
	AgeX  int    `src:"Age"`
	AgeY  int    `rule:"age" src:"Age"`
	InfoX string `rule:"info" src:"Name,Info[0],Info[1]"`
}

func getAge(age int) int {
	return age + 1
}

func getInfo(a, b, c string) string {
	return a + b + c
}

func TestTransfer(t *testing.T) {
	RegisterRule("age", getAge)
	RegisterRule("info", getInfo)
	err := CheckRules()
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	src := &Src{
		Name: "雪辙",
		Age:  18,
		Info: []string{"a", "b"},
	}

	dst := &Dst{}

	err = Transfer(src, dst)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	t.Logf("dst: %+v", dst)
}
