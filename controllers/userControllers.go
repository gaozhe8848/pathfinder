package controllers

import (
	"net/http"
	"os"
	"pathfinder/initializers"
	"pathfinder/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "hello!"})
}

func SignUp(ctx *gin.Context) {

	var reqBody struct {
		Email    string
		Password string
	}
	//bind req to reqBody
	if ctx.Bind(&reqBody) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed read body",
		})
		return
	}

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	//create the user
	user := models.User{Email: reqBody.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	//respond
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created!",
	})

}

func Login(ctx *gin.Context) {

	var user models.User

	var reqBody struct {
		Email    string
		Password string
	}

	//bind req to reqBody
	if ctx.Bind(&reqBody) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed read body",
		})
		return
	}

	//look up for the user
	initializers.DB.First(&user, "email = ?", reqBody.Email)

	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid username or pwd",
		})
		return
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid username or pwd",
		})
		return
	}

	//generate a JWT token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create JWT token",
		})
		return
	}

	//respond with JWT as a cookie

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{})

}

func Validate(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
