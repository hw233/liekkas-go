package strings

import "testing"

func TestStringDisplayLen(t *testing.T) {
	t.Log(StringDisplayLen("不会吧123213"))
	t.Log(StringDisplayLen("123213"))
	t.Log(StringDisplayLen("ABC"))
	t.Log(StringDisplayLen("abc"))
	t.Log(StringDisplayLen("。。。"))
	t.Log(StringDisplayLen("..."))

}
