package core

import (
	"context"

	"net/http"

	"strings"

	"time"

	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
)

/*
TokenClaims holds information for token
*/
type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

/*
CreateToken creates auth token
*/
func CreateToken(db *gorm.DB, user User, secret string, expiration int) (result string, err error) {

	claims := TokenClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			// what to do with this
			ExpiresAt: time.Now().Unix() + int64(expiration),
			Issuer:    TOKEN_ISSUER,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err = token.SignedString([]byte(secret))

	return
}

/*
ParseToken parses token and returns claims
*/
func ParseToken(r *http.Request, secret string) (claims *TokenClaims, err error) {

	var (
		token       *jwt.Token
		tokenString string
	)

	if tokenString, err = GetRequestToken(r); err != nil {
		return
	}

	token, err = jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// check errors and return appropriate error
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = ErrTokenInvalid
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				err = ErrTokenExpired
				return
			} else {
				err = ErrTokenInvalid
				return
			}
		} else {
			err = ErrTokenInvalid
			return
		}
	}

	// set claims
	claims = token.Claims.(*TokenClaims)

	return
}

/*
GetRequestToken returns token from request
*/
func GetRequestToken(r *http.Request) (result string, err error) {
	prefix := "Bearer "
	header := r.Header.Get(TOKEN_HEADER_NAME)
	if !strings.HasPrefix(header, prefix) {
		err = ErrTokenInvalid
		return
	}
	result = strings.TrimSpace(header[len(prefix):])

	return
}

/*
Return token claims from request context
*/
func ContextGetTokenUser(ctx context.Context) (user User, err error) {
	result := ctx.Value(CONTEXT_TOKEN_USER)

	var ok bool
	if user, ok = result.(User); !ok {
		err = ErrTokenUserInvalid
	}

	return
}

/*
Return token claims from request context
*/
func ContextSetTokenUser(ctx context.Context, user User) (result context.Context) {
	return context.WithValue(ctx, CONTEXT_TOKEN_USER, user)
}
