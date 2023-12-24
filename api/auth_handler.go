package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore   db.UserStore
	doctorStore db.DoctorStore
}

func NewAuthHandler(userStore db.UserStore, doctorStore db.DoctorStore) *AuthHandler {
	return &AuthHandler{
		userStore:   userStore,
		doctorStore: doctorStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}
type DoctorAuthResponse struct {
	Doctor *types.Doctor `json:"doctor"`
	Token  string        `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func (h *AuthHandler) HandleRegisterUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, _ := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if user != nil {
		return c.JSON(fiber.Map{"error": "user already present"})
	}
	newuser, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), newuser)
	if err != nil {
		return invalidCredentials(c)
	}
	return c.JSON(insertedUser)
}
func (h *AuthHandler) HandleRegisterDoctor(c *fiber.Ctx) error {
	var params types.CreateDoctorParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.ValidateDoctor(); len(errors) > 0 {
		return c.JSON(errors)
	}
	doctor, _ := h.doctorStore.GetDoctorByEmail(c.Context(), params.Email)
	if doctor != nil {
		return c.JSON(fiber.Map{"error": " already exists"})
	}
	newdoctor, err := types.NewDoctorFromParams(params)
	fmt.Println(newdoctor.ID.Hex())
	if err != nil {
		return c.JSON(fiber.Map{"error at creating new doctor from params :": err.Error()})
	}
	insertedDoctor, err := h.doctorStore.InsertDoctor(c.Context(), newdoctor)
	if err != nil {
		return c.JSON(fiber.Map{"error at invalid user insertion:": err.Error()})
	}
	return c.JSON(insertedDoctor)
}

func (h *AuthHandler) HandleLoginUser(c *fiber.Ctx) error {
	cookievalue := c.Cookies("jwt")
	if cookievalue != "" {
		fmt.Println(cookievalue)
		return c.JSON(fiber.Map{"error": "token exists in cookie logout to continue"})
	}
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}
	resp := UserAuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    resp.Token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(resp)
}

func (h *AuthHandler) HandleLoginDoctor(c *fiber.Ctx) error {
	cookievalue := c.Cookies("jwt")
	if cookievalue != "" {
		fmt.Println(cookievalue)
		return c.JSON(fiber.Map{"error": "token exists in cookie logout to continue"})
	}
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	doctor, err := h.doctorStore.GetDoctorByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}
	if !types.IsValidPassword(doctor.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}
	resp := DoctorAuthResponse{
		Doctor: doctor,
		Token:  CreateTokenFromDoctor(doctor),
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    resp.Token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(resp)
}

func (h *AuthHandler) HandleUserLogOut(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  expired,
		HTTPOnly: true,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *AuthHandler) HandleDoctorLogOut(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  expired,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 6).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}

func CreateTokenFromDoctor(doctor *types.Doctor) string {
	now := time.Now()
	expires := now.Add(time.Hour * 6).Unix()
	claims := jwt.MapClaims{
		"id":      doctor.ID.Hex(),
		"email":   doctor.Email,
		"expires": expires,
	}
	fmt.Println(claims["id"])
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
