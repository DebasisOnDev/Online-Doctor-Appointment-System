package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DoctorID          primitive.ObjectID `bson:"doctorId" json:"doctorId"`
	UserID            primitive.ObjectID `bson:"userId" json:"userId"`
	Department        string             `bson:"Department" json:"Department"`
	Issue             string             `bson:"issue" json:"issue"`
	Date              time.Time          `bson:"date" json:"date"`
	AppointmentTiming string             `bson:"appTiming" json:"appTiming"`
	AppointmentNumber int                `bson:"appNumber" json:"appNumber"`
	Fee               int                `bson:"fee" json:"fee"`
	IsPaid            bool               `bson:"isPaid" json:"isPaid"`
	IsComplete        bool               `bson:"isComplete" json:"isComplete"`
}
