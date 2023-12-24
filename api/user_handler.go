package api

import (
	"errors"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (u *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := u.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrNotResourceNotFound("user")
	}
	return c.JSON(users)
}

func (u *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := u.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(fiber.Map{"error": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (u *UserHandler) HandleInsertUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := u.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
