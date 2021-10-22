package models

import (
	"errors"
	"strings"
	"time"

	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/validators"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

const (
	ACCESS_SECRET  = "jdnfksdmfksd"
	REFRESH_SECRET = "mcmvmkmsdnfsdmfdsjf"
)

func ExtractToken(header string) string {
	strArr := strings.Split(header, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(token string) (*jwt.Token, error) {
	tokenString := ExtractToken(token)
	tokenStruct, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("undexpected signing alg")
		}

		return []byte(ACCESS_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	return tokenStruct, nil
}

func TokenValid(header string) error {
	token, err := VerifyToken(header)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(header string) (*validators.AccessDetails, error) {
	token, err := VerifyToken(header)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userUuid, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}

		return &validators.AccessDetails{
			AccessUuid: accessUuid,
			UserUuid:   userUuid,
		}, nil
	}
	return nil, err
}

func CreateToken(ID string) (*validators.TokenDetails, error) {
	td := &validators.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 60).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = ID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(ACCESS_SECRET))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = ID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(REFRESH_SECRET))
	if err != nil {
		return nil, err
	}
	err = db.SaveToken(ID, td)
	if err != nil {
		return nil, err
	}

	return td, nil
}
