package auth

import (
	"errors"
	"fmt"
	"gin_example/core"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Sigin(ctx *gin.Context) {
	conn := core.GetDatabase()
	defer core.CloseDB(conn)
	var tuser struct {
		Name     string
		Email    string
		Password string
	}

	tuser.Email = ctx.PostForm("Email")
	tuser.Name = ctx.PostForm("Name")
	tuser.Password = ctx.PostForm("Password")

	if len(tuser.Name) <= 8 && len(tuser.Password) <= 8 && !isEmailValid(tuser.Email) {
		err := "Eror in request"
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}

	var dbuser core.User
	conn.Where("enail = ?", tuser.Email).First(&dbuser)

	err := conn.Where("enail = ?", tuser.Email).First(&dbuser).Error
	errors.As(err, gorm.ErrRecordNotFound)
	if dbuser.Email != "" {
		err := "Eror in request"
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}

	check := checkPasswd(tuser.Password, dbuser.Password)
	if !check {
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": "Error password",
		})
	}

	validtoken, err := GenerateJWT(dbuser.Email, "admin")

	var token core.Token
	token.Email = dbuser.Email
	token.Role = "admin"
	token.TokenString = validtoken

	ctx.JSON(http.StatusOK, gin.H{"token": token})

}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func checkPasswd(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}

func GenerateJWT(email, rol string) (string, error) {
	var mysalt = []byte("Hola_mundo")
	token := jwt.New(jwt.SigningMethodES256)
	claim := token.Claims.(jwt.MapClaims)

	claim["authorized"] = true
	claim["email"] = email
	claim["role"] = rol
	claim["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenstring, err := token.SignedString(mysalt)
	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenstring, nil
}
