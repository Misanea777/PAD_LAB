package validators

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	val := validator.New()
	err := val.RegisterValidation("passwd", func(f1 validator.FieldLevel) bool {
		return len(f1.Field().String()) > 6
	})
	if err != nil {
		panic(err.Error())
	}
	err = val.RegisterValidation("usernme", func(f1 validator.FieldLevel) bool {
		return len(f1.Field().String()) > 9
	})
	if err != nil {
		panic(err.Error())
	}
	Validate = val
}

type UserCredentials struct {
	Username string `json:"username" validate:"usernme"`
	Password string `json:"password" validate:"passwd"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

type AccessDetails struct {
	AccessUuid string
	UserUuid   string
}

type UserBalanceInfo struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
}
