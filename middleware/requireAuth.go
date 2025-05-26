package middleware

import (
	"fmt"
	"log"
	"net/http"
	"pathfinder/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func RequireAuth(db *gorm.DB, jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//get the cookie off req
		tokenString, err := ctx.Cookie("Authorization")

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//decode/validate
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})
		if err != nil {
			log.Print(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			//check the exp
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				log.Println("jwt within expiry") // This log message seems to indicate the opposite of the check
				return
			}

			//find the user with the token sub
			var user models.User
			db.First(&user, claims["sub"]) // Use injected db

			if user.ID == 0 {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				log.Println("failed to find user")
				return
			}

			//attach to req
			ctx.Set("user", user)

			//continue
			ctx.Next()
		} else {
			log.Println("Invalid JWT claims")
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
