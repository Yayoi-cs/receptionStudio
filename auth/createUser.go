package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func CreateUser(email string) string {

	accessToken, err := generateJwtAuthStep1(email)
	if err != nil {
		return ""
	}
	return accessToken

}

func GenerateJwtAuthGeneral(email string) (string, error) {
	expirationTime := time.Now().AddDate(0, 1, 0)
	claims := jwt.MapClaims{
		"username": email,
		"exp":      expirationTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := jwtSecret()
	accessToken, err := token.SignedString(secret)
	return accessToken, err
}

func generateJwtAuthStep1(email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := jwt.MapClaims{
		"username": email,
		"exp":      expirationTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := jwtSecret()
	accessToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return accessToken, err
}

func CheckJwt(tokenString string) (string, error) {
	jwtSecret := jwtSecret()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("JWT is Valid")
	}
	expDate, ok := claims["exp"].(string)
	if !ok {
		return "", fmt.Errorf("can't get expirationDate")
	}
	expTime, err := time.Parse(time.RFC3339Nano, expDate)
	if err != nil {
		return "", fmt.Errorf("can't parse expirationDate")
	} else {
		now := time.Now()
		if now.After(expTime) {
			return "", fmt.Errorf("JWT is expired")
		}
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("can't get Username from jwt")
	}

	return username, nil
}

func SendConfirmEmail(address string, confirmCode string) {
	from := "hirogoshawk3249@gmail.com"
	subject := "Hello"
	body := "Hello World!\nYour verify code :" + confirmCode
	hostname, port, username, password := config()
	auth := smtp.PlainAuth("", username, password, hostname)
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", address, subject, body), "\n", "\r\n"))
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", hostname, port), auth, from, []string{address}, msg); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
