package uuid

import (
	uuid "github.com/satori/go.uuid"
)

var prefixUUID uuid.UUID

func init() {
	prefixUUID = uuid.NewV1()
}

func GenUUID() string {
	return uuid.NewV5(prefixUUID, uuid.NewV4().String()).String()
}
