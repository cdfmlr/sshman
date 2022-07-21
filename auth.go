package main

// This file contains the implementation of a JWT authenticator.
// And we will use it in router.go to verify the user's JWT token.

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cdfmlr/crud/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Payload is the JWT payload (i.e. information store in the token).
type Payload struct {
	Id    uint   `json:"id"`
	Email string `json:"email"`
	Role  int    `json:"role"`
}

// Claims is a JWT claims with Payload.
// HEADER:
//    { "alg": "HS256", "typ": "JWT" }
// PAYLOAD:
//    {
//	    "id": 1,
//	    "email": "root@sshman.example",
//	    "role": 3,
//	    "exp": 1658372957
//    }
// SIGNATURE:
//    ...
type Claims struct {
	*jwt.RegisteredClaims
	Payload
}

//////////// JwtAuthenticator  ////////////
// u may seek to different authentication
// procedures, such as RSA public key authentication,
// so we will use an interface to abstract it.

// JwtAuthenticator verifies user's JWT token and parses the payload.
type JwtAuthenticator interface {
	AuthAndParseToken(token string) (*Payload, error)
}

//////////// HMAC Sample Secret Auth ////////////

// HmacSecretAuthenticator is an implementation of JwtAuthenticator.
// From: https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac
type HmacSecretAuthenticator struct {
	hmacSampleSecret []byte
}

func NewHmacSecretAuthenticator(secret string) JwtAuthenticator {
	return &HmacSecretAuthenticator{
		hmacSampleSecret: []byte(secret),
	}
}

func (a *HmacSecretAuthenticator) AuthAndParseToken(token string) (*Payload, error) {
	logger := log.ZoneLogger("sshman/auth")

	token = strings.TrimSpace(token)

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return a.hmacSampleSecret, nil
	})
	if err != nil {
		logger.Errorf("AuthAndParseToken error: %v", err)
		return nil, err
	}
	return &claims.Payload, nil
}

//////////// Middlewares ////////////

// JwtAuthMiddleware verifies user's JWT token.
func JwtAuthMiddleware(authenticator JwtAuthenticator) gin.HandlerFunc {
	logger := log.ZoneLogger("sshman/auth")

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		userInfo, err := authenticator.AuthAndParseToken(token)
		if err != nil {
			logger.Errorf("JwtAuthMiddleware Handler error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set("userInfo", userInfo)
		c.Next()
	}
}

// RoleAuthMiddleware checks if the user has the required role permission.
func RoleAuthMiddleware(role int) gin.HandlerFunc {
	logger := log.ZoneLogger("sshman/auth")
	return func(c *gin.Context) {
		userInfo := c.MustGet("userInfo").(*Payload)
		if userInfo.Role&role == 0 {
			logger.WithField("userInfo", userInfo).
				WithField("role", role).
				Errorf("RoleAuthMiddleware auth failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		c.Next()
	}
}
