package auth

import (
	"aidanwoods.dev/go-paseto"
)

type PasetoManager struct {
	key paseto.V4SymmetricKey
}

func NewPasetoManager(key paseto.V4SymmetricKey) *PasetoManager {
	return &PasetoManager{
		key: key,
	}
}
