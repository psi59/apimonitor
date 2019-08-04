package rsstr

import (
	"strings"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}

func NewUUIDWithoutHyphen() string {
	return strings.Replace(NewUUID(), "-", "", -1)
}
