package uid

import "testing"

func TestUID(t *testing.T) {
	uid := NewUID()
	t.Logf("uid: %d", uid.Gen())
	t.Logf("uid: %d", uid.Gen())
}
