package whitelist

import "testing"

func TestMultiWhiteList_Filter(t *testing.T) {
	whiteList := NewMultiWhiteList()
	ops := NewOpsWithUid(1)
	t.Log(whiteList.Filter(ops))

	whiteList.Reload(ops)
	t.Log(whiteList.Filter(ops))
	whiteList.Del(ops)
	t.Log(whiteList.Filter(ops))
	whiteList.Add(ops)
	t.Log(whiteList.Filter(ops))

}
