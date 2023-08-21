package origin

// this package manage origin platform account, Pylon

import (
	"github.com/PylonSchema/server/auth"
	"github.com/PylonSchema/server/model"
	"github.com/google/uuid"
)

type Database interface {
	CreateOriginUser(user *model.User, origin *model.Origin) error
	IsEmailUsed(email string) (bool, error)
	GetOriginUser(email string, password string) (*model.User, error)
}

type AuthOriginAPI struct {
	DB      Database
	JwtAuth *auth.JwtAuth
}

type loginPaylaod struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func New(db Database, jwtAuth *auth.JwtAuth) *AuthOriginAPI {
	return &AuthOriginAPI{
		DB:      db,
		JwtAuth: jwtAuth,
	}
}

func (c *createPayload) createModel(hashedPassword string) (*model.User, *model.Origin, error) {
	userUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, err
	}
	return &model.User{
			Username:    c.Username,
			AccountType: 1,
			UUID:        userUUID,
			Email:       c.Email,
		}, &model.Origin{
			UUID:     userUUID,
			Password: hashedPassword,
		}, nil
}
