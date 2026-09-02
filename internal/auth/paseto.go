package auth

import (
	"fmt"
	"strconv"
	"time"

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

func (p *PasetoManager) CreateToken(userID int) (string, error) {
	token := paseto.NewToken()

	token.SetString("user_id", fmt.Sprintf("%d", userID))
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(24 * time.Hour))

	return token.V4Encrypt(p.key, nil), nil
}

func (p *PasetoManager) VerifyToken(token string) (int, error) {
	parsedToken, err := paseto.NewParser().ParseV4Local(p.key, token, nil)
	if err != nil {
		return 0, nil
	}

	userStringID, err := parsedToken.GetString("user_id")
	if err != nil {
		return 0, nil
	}

	userID, err := strconv.Atoi(userStringID)
	if err != nil {
		return 0, nil
	}

	return userID, nil
}
