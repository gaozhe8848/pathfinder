package controllers

import (
	"errors"
	"net/http"
	"pathfinder/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn" // Or your specific DB driver's error type for unique constraints
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthController holds the database dependency.
type AuthController struct {
	DB        *gorm.DB
	JWTSecret string
}

// NewAuthController creates a new AuthController with the given database connection.
func NewAuthController(db *gorm.DB, jwtSecret string) *AuthController {
	return &AuthController{DB: db, JWTSecret: jwtSecret}
}

func (ac *AuthController) SignUp(ctx *gin.Context) {

	var reqBody struct {
		Email    string
		Password string
	}
	//bind req to reqBody
	if ctx.Bind(&reqBody) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to process password",
		})
		return
	}

	//create the user
	user := models.User{Email: reqBody.Email, Password: string(hash)}

	result := ac.DB.Create(&user)

	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation for PostgreSQL
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
		}
		return
	}

	//respond
	ctx.JSON(http.StatusCreated, gin.H{
		// Consider returning the created user resource or a success message
		"message": "User created successfully",
	})

}

func (ac *AuthController) Login(ctx *gin.Context) {

	var user models.User

	var reqBody struct {
		Email    string
		Password string
	}

	//bind req to reqBody
	if ctx.Bind(&reqBody) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	//look up for the user
	result := ac.DB.First(&user, "email = ?", reqBody.Email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email or password",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error while trying to find user",
			})
		}
		return
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	//generate a JWT token
	secret := ac.JWTSecret
	if secret == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured on server"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create JWT token",
		})
		return
	}

	//respond with JWT as a cookie
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24, "/", "", false, true) // Secure flag is now always false
	ctx.JSON(http.StatusOK, gin.H{})

}

func (ac *AuthController) Validate(ctx *gin.Context) {
	userValue, exists := ctx.Get("user")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated or session expired",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": userValue, // Return the user object under "user" key
	})
}
