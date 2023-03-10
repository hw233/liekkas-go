package global

import "shared/utility/key"

const (
	globalCacheKeyPrefix  = "global:cache"
	globalHashKeyPrefix   = "global:hash"
	globalIncrIDKeyPrefix = "global:id"
	globalSetKeyPrefix    = "global:set"
	globalLockKeyPrefix   = "global:lock"
	globalStringKeyPrefix = "global:string"
	globalListKeyPrefix   = "global:list"
)

func makeSetKey(key interface{}) string {
	return makeKey(globalSetKeyPrefix, key)
}

func makeHashKey(key interface{}) string {
	return makeKey(globalHashKeyPrefix, key)
}

func makeLockKey(key interface{}) string {
	return makeKey(globalLockKeyPrefix, key)
}

func makeStringKey(key interface{}) string {
	return makeKey(globalStringKeyPrefix, key)
}

func makeIncrIDKey(key interface{}) string {
	return makeKey(globalIncrIDKeyPrefix, key)
}

func makeListKey(key interface{}) string {
	return makeKey(globalListKeyPrefix, key)
}

func makeKey(vals ...interface{}) string {
	return key.MakeRedisKey(vals...)
}
