package models

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	ID       uint
	NickName string
	// 权限id
	AuthorityId uint
	jwt.StandardClaims
}
