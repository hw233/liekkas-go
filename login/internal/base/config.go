package base

const (
	ServerENVTest = "test"
)

func IsTestServerENV() bool {
	return Config.ServerENV == ServerENVTest
}
