package auth

import (
	"aidanwoods.dev/go-paseto"
)

type PasetoManager struct {
	key paseto.V4SymmetricKey
}
