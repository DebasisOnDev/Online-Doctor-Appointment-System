package api

import (
	"fmt"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type CheckUpHandler struct {
	store *db.Store
}

func NewCheckUpHandler(store *db.Store) *CheckUpHandler {
	return &CheckUpHandler{
		store: store,
	}
}

func (ch *CheckUpHandler) HandleGetAllCheckups(c *fiber.Ctx) error {
	dr, err := getAuthDoctor(c)
	if err != nil {
		return c.JSON(fiber.Map{"error": "unauthenticated"})
	}
	checkups, err := ch.store.CheckUp.GetCheckUps(c.Context(), dr.ID.String())
	if err != nil {
		fmt.Println("error at getting checkups")
		return err
	}
	return c.JSON(checkups)
}

func (ch *CheckUpHandler) HandleDoCheckUp(c *fiber.Ctx) error {
	id := c.Params("id")
	bookings, err := ch.store.Booking.GetBookings(c.Context(), id)
	if err != nil {
		return err
	}
	var ncheckup []*types.CheckUp
	for _, book := range bookings {
		checkup, err := ch.store.CheckUp.GetCheckByID(c.Context(), book.ID.Hex())
		if err != nil {
			return err
		}
		filter := bson.D{{Key: "_id", Value: checkup.ID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: true}, {Key: "fee", Value: book.Fee}}}}
		cu, err := ch.store.CheckUp.PerformCheckUp(c.Context(), filter, update)
		if err != nil {
			return err
		}
		ncheckup = append(ncheckup, cu)
		if err := ch.store.Booking.UpdateBooking(c.Context(), book.ID.Hex()); err != nil {
			return err
		}
	}
	return c.JSON(ncheckup)

}

func (ch *CheckUpHandler) HandleGetCheckUpByID(c *fiber.Ctx) error {
	id := c.Params("id")
	checkups, err := ch.store.CheckUp.GetCheckByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(checkups)
}
