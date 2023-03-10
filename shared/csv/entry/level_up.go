package entry

type LevelCfgCache interface {
	//GetExpToNextLevel(nowLevel int32) (int32, bool)

	GetExpArr() []int32
}
