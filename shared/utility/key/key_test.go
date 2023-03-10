package key

import "testing"

func TestMakeKey(t *testing.T) {
	t.Logf(MakeRedisKey("key1", "key2"))
	t.Logf(MakeEtcdKey("key1"))
}
