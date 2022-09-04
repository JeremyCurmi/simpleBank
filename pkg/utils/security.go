package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/JeremyCurmi/simpleBank/pkg/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TokenValid(c *gin.Context) error {
	_, err := ParseJWTToken(c)
	return err
}

func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"exp":        time.Now().Add(time.Hour * time.Duration(*config.TokenHourLifeSpan)).Unix(),
	})
	return token.SignedString([]byte(*config.APISecretKey))
}

// ParseJWTToken checks that token is of correct format
func ParseJWTToken(c *gin.Context) (*jwt.Token, error) {
	token := ExtractToken(c)

	jwtKeyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(*config.APISecretKey), nil
	}
	return jwt.Parse(token, jwtKeyFunc)
}

// ExtractToken tries to extract bearer token from request
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID extracts user id from request
func ExtractTokenID(c *gin.Context) (uint, error) {
	token, err := ParseJWTToken(c)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
