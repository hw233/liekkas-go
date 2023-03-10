package base

type Config interface {
	// Load file
	Load(path string) bool

	// clear data
	Clear()
}
