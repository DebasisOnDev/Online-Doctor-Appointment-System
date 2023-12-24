package api

import (
	"errors"
	"fmt"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DoctorHandler struct {
	doctorStore  db.DoctorStore
	checkupStore db.CheckUpStore
}

func NewDoctorHandler(doctorStore db.DoctorStore, checkupStore db.CheckUpStore) *DoctorHandler {
	return &DoctorHandler{
		doctorStore:  doctorStore,
		checkupStore: checkupStore,
	}
}

func (d *DoctorHandler) GetCheckups(c *fiber.Ctx) error {
	dr, err := getAuthDoctor(c)
	if err != nil {
		return err
	}
	ch, err := d.checkupStore.GetCheckUps(c.Context(), dr.ID.String())
	if err != nil {
		fmt.Println("fetch checkup error")
		return err
	}
	return c.JSON(ch)
}

func (d *DoctorHandler) HandleGetDoctors(c *fiber.Ctx) error {
	dr, err := getAuthDoctor(c)
	fmt.Println(dr)
	if err != nil {
		return err
	}
	filter := bson.M{"specialist": dr.Specialist}
	doctors, err := d.doctorStore.GetDoctorsBySpecialist(c.Context(), filter)
	if err != nil {
		return ErrNotResourceNotFound("doctor")
	}
	return c.JSON(doctors)
}

func (d *DoctorHandler) HandleGetAllDoctorsByUser(c *fiber.Ctx) error {
	doctors, err := d.doctorStore.GetDoctors(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(doctors)
}

func (d *DoctorHandler) HandleGetDoctor(c *fiber.Ctx) error {
	id := c.Params("id")
	doctor, err := d.doctorStore.GetDoctorByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(fiber.Map{"error": "not found"})
		}
		return err
	}
	return c.JSON(doctor)
}

func (d *DoctorHandler) HandleInsertDoctor(c *fiber.Ctx) error {
	var params types.CreateDoctorParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.ValidateDoctor(); len(errors) > 0 {
		return c.JSON(errors)
	}
	doctor, err := types.NewDoctorFromParams(params)
	if err != nil {
		return err
	}
	insertedDoctor, err := d.doctorStore.InsertDoctor(c.Context(), doctor)
	if err != nil {
		return err
	}
	return c.JSON(insertedDoctor)
}
