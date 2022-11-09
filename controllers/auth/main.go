package auth

import (
	"errors"
	"fmt"
	"gin_example/core"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SigUp(ctx *gin.Context) {
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
	errors.Is(err, gorm.ErrRecordNotFound)

	if dbuser.Email != "" {
		err := "Eror in request"
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}

	tuser.Password, err = generatehashPassword(tuser.Password)
	if err != nil {
		log.Fatal("error in password hash")
	}

	var user core.User = core.User{
		Model:         gorm.Model{},
		Name:          tuser.Name,
		Email:         tuser.Email,
		Password:      tuser.Password,
		Role:          "Default",
		LastConection: ctx.RemoteIP(),
		TimeSession:   strconv.FormatInt(time.Now().Unix(), 10),
	}
	conn.Create(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"Usuer": &user,
	})

}

func SigIn(ctx *gin.Context) {
	conn := core.GetDatabase()
	defer core.CloseDB(conn)

	var aute core.Authentication

	aute.Email = ctx.PostForm("Email")
	aute.Password = ctx.PostForm("Password")

	if len(aute.Password) <= 8 && !isEmailValid(aute.Email) {
		err := "Empty var"
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}

	var authuser core.User
	conn.Where("Email = ?", aute.Email).First(&authuser)
	if authuser.Email == "" {
		err := "Empty Email"
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}

	check := checkPasswd(aute.Password, authuser.Password)
	fmt.Printf("Error in pass :%v \n", check)
	fmt.Printf("Pass :%v \n", aute.Password)
	if !check {
		err := "Empty Password"
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}

	validtoken, err := GenerateJWT(authuser.Email, "admin")
	if err != nil {
		err := "Error generate Token"
		ctx.JSON(http.StatusNotFound, gin.H{
			"err": err,
		})
		return
	}
	var token core.Token
	token.Email = authuser.Email
	token.Role = authuser.Name
	token.TokenString = validtoken
	ctx.JSON(http.StatusOK, gin.H{
		"Token": token,
	})
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func checkPasswd(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}

func generatehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateJWT(email, role string) (string, error) {
	var mysalt = []byte("Hola_mundo")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mysalt)

	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Auth(ctx *gin.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if len(ctx.Request.Header.Get("Token")) > 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "No Token Found",
			})
			return
		}
		var mysalt = []byte("Hola_mundo")

		token, err := jwt.Parse(ctx.Request.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				return nil, fmt.Errorf("Error parsin")
			}
			return mysalt, nil
		})

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Token has been expired",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "Admin" {
				ctx.Request.Response.Header.Add("Role", "Admin")
				return

			} else if claims["role"] == "Default" {
				ctx.Request.Response.Header.Add("Role", "user")
				return
			}
		}

		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Not Authorized",
		})

	}
}
