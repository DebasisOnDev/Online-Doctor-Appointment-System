package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookingStore interface {
	GetAllDoctorsByUser(context.Context, bson.M) ([]*types.Doctor, error)
	GetADoctorByUser(context.Context, string) (*types.Doctor, error)
	BookADoctorByUser(context.Context, *types.Doctor, *types.User) (*types.Booking, error)
	//GetDoctorBySpecialist(context.Context, string) ([]*types.Doctor, error)
	GetAllDoctor(context.Context) ([]*types.Doctor, error)
	GetBookings(context.Context, string) ([]*types.Booking, error)
	UpdateBooking(context.Context, string) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	DoctorStore
	CheckUpStore
	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoBookingStore{
		client:       client,
		coll:         client.Database(dbname).Collection("bookings"),
		CheckUpStore: NewMongoCheckUpStore(client),
	}
}

func (m *MongoBookingStore) GetBookings(ctx context.Context, id string) ([]*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": oid}
	cur, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (m *MongoBookingStore) GetAllDoctorsByUser(ctx context.Context, filter bson.M) ([]*types.Doctor, error) {
	doc, err := m.DoctorStore.GetDoctorsBySpecialist(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (m *MongoBookingStore) UpdateBooking(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	var booking types.Booking
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "Fee", Value: 500},
		{Key: "isPaid", Value: true}, {Key: "isComplete", Value: true},
	}}}
	err = m.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&booking)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoBookingStore) BookADoctorByUser(ctx context.Context, doctor *types.Doctor, user *types.User) (*types.Booking, error) {
	booking := &types.Booking{}
	booking.DoctorID = doctor.ID
	booking.UserID = user.ID
	booking.Department = doctor.Specialist
	booking.Date = time.Now().Add(time.Hour * 24)
	booking.AppointmentTiming = doctor.WorkingHour
	booking.AppointmentNumber = 1
	booking.Fee = doctor.Fee
	booking.IsPaid = false
	booking.IsComplete = false

	book, err := m.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, fmt.Errorf("failed to insert booking %v", err)
	}

	booking.ID = book.InsertedID.(primitive.ObjectID)

	if m.CheckUpStore == nil {
		return nil, fmt.Errorf("m.Store.CheckUp is nil")
	}
	checkup, err := types.NewCheckUpFromBooking(booking)
	if err != nil {
		return nil, err
	}

	//ch, err := m.CheckUp.InsertCheckUp(ctx, checkup)
	_, err = m.CheckUpStore.InsertCheckUp(ctx, checkup)

	if err != nil {

		return nil, err
	}

	return booking, nil

}

func (m *MongoBookingStore) GetADoctorByUser(ctx context.Context, id string) (*types.Doctor, error) {
	dr, err := m.DoctorStore.GetDoctorByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dr, nil
}

func (m *MongoBookingStore) GetAllDoctor(ctx context.Context) ([]*types.Doctor, error) {
	fmt.Println("error 1")
	docs, err := m.DoctorStore.GetDoctors(ctx)
	fmt.Println("error 2")
	if err != nil {
		return nil, err
	}
	return docs, nil
}
