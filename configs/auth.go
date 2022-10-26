package configs

import (
	jwtware "github.com/gofiber/jwt/v3"
)

var secretKey = Config("SECRET_KEY")

var ConfigAuth = jwtware.New(jwtware.Config{
	SigningKey: []byte(secretKey),
})
