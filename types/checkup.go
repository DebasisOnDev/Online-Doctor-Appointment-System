package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type CheckUp struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DoctorID primitive.ObjectID `bson:"doctorId" json:"doctorId"`
	UserID   primitive.ObjectID `bson:"userId" json:"userId"`
	Status   bool               `bson:"status" json:"status"`
	Fee      float64            `bson:"fee" json:"fee"`
}

func NewCheckUpFromBooking(booking *Booking) (*CheckUp, error) {
	return &CheckUp{
		DoctorID: booking.DoctorID,
		UserID:   booking.UserID,
		Status:   true,
		Fee:      500,
	}, nil
}
