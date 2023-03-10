package mysql

type Marshaler interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
