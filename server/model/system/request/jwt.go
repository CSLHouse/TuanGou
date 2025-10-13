package request

import (
	"github.com/gofrs/uuid/v5"
	jwt "github.com/golang-jwt/jwt/v4"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          int
	UserName    string
	AuthorityId int
	Telephone   string
	CaptainId   int
}

type WXBaseClaims struct {
	OpenId   string
	NickName string
}

type WXAccountClaims struct {
	WXBaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}
