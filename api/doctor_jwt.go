package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/gofiber/fiber/v2"
)

func JWTAuthenticationDoctor(doctorStore db.DoctorStore) fiber.Handler {
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
		fmt.Println(claims)
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		doctorID, ok := claims["id"].(string)
		fmt.Println(doctorID)
		if !ok {
			fmt.Println("Error: ID is not a string")
		}
		doctor, err := doctorStore.GetDoctorByID(c.Context(), doctorID)
		if err != nil {
			fmt.Println("error at setting token")
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("doctor", doctor)
		return c.Next()
	}
}
