package ip

import "testing"

func TestListenAddr(t *testing.T) {
	addr, err := ListenAddr()
	t.Logf("%s, %v", addr, err)
}
