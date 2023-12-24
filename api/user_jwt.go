package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// func JwtAuth(userStore db.UserStore) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		cookie := c.Cookies("jwt")
// 		sec := os.Getenv("JWT_SECRET")
// 		log.Panic("entered the authentication", cookie, sec)
// 		token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(sec), nil
// 		})
// 		if err != nil {
// 			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"message": "unauthorized",
// 			})
// 		}
// 		claims := token.Claims.(*jwt.MapClaims)
// 		//fmt.Println(claims, "failed to get issuer")
// 		_, err = claims.GetIssuer()
// 		if err != nil {
// 			//fmt.Println(err.Error())
// 			return err
// 		}
// 		//fmt.Println(st)
// 		return c.Next()

// 	}
// }

func JWTAuthenticationUser(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {

		cookie := c.Cookies("jwt")
		fmt.Println(cookie)

		if cookie == "" {
			fmt.Println("cookie is empty")
			return ErrUnAuthorized()
		}

		claims, err := validateToken(cookie)
		if err != nil {
			fmt.Println("invalid token")
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			fmt.Println("invalid user found ")
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token", err)
		return nil, ErrUnAuthorized()
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, ErrUnAuthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnAuthorized()
	}
	return claims, nil
}
