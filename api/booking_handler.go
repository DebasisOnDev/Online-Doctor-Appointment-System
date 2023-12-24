package api

import (
	"fmt"
	"time"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (b *BookingHandler) HandleGetDoctorsByUser(c *fiber.Ctx) error {
	var doctors []*types.Doctor
	doctors, err := b.store.Doctor.GetDoctors(c.Context())
	if err != nil {
		return ErrNotResourceNotFound("bookings")
	}

	return c.JSON(doctors)
}

func (b *BookingHandler) HandleGetDoctorBySpecialist(c *fiber.Ctx) error {
	specialist := c.Params("specialist")
	sp := bson.M{"specialist": specialist}

	doctors, err := b.store.Doctor.GetDoctorsBySpecialist(c.Context(), sp)
	if err != nil {
		return ErrNotResourceNotFound("bookings")
	}
	return c.JSON(doctors)
}

func (b *BookingHandler) HandleGetDoctorById(c *fiber.Ctx) error {
	id := c.Params("id")
	doctor, err := b.store.Doctor.GetDoctorByID(c.Context(), id)
	if err != nil {
		return c.JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(doctor)
}

func (b *BookingHandler) HandleBookingADoctorByUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user *types.User
	var doctor *types.Doctor

	doctor, err := b.store.Doctor.GetDoctorByID(c.Context(), id)
	if err != nil {

		return ErrUnAuthorized()
	}

	user, err = getAuthUser(c)

	if err != nil {

		return ErrUnAuthorized()
	}

	if user == nil || doctor == nil {

		return ErrUnAuthorized()
	}

	var booking *types.Booking
	booking, err = b.store.Booking.BookADoctorByUser(c.Context(), doctor, user)

	if err != nil {

		return ErrUnAuthorized()
	}

	doctor.AppointmentInfo.AppointmentNumber = booking.AppointmentNumber
	doctor.AppointmentInfo.AppointmentDate = time.Now().Add(time.Hour * 24)
	doctor.AppointmentInfo.Patient = user.ID.String()

	if dr := b.store.Doctor.SetDoctorAppointmentInfo(c.Context(), doctor, user.ID.String()); dr == nil {
		fmt.Println("no doctor fetched")
	}
	return c.JSON(booking)
}

func (b *BookingHandler) HandleGetAllDoctors(c *fiber.Ctx) error {
	docs, err := b.store.Booking.GetAllDoctor(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(docs)
}
