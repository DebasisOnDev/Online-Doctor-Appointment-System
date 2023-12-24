package api

import (
	"fmt"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return user, nil
}

func getAuthDoctor(c *fiber.Ctx) (*types.Doctor, error) {
	doctor, ok := c.Context().UserValue("doctor").(*types.Doctor)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return doctor, nil
}
