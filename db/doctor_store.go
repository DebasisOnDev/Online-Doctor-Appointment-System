package db

import (
	"context"
	"fmt"
	"os"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const doctorColl = "doctors"

type DoctorStore interface {
	GetDoctors(context.Context) ([]*types.Doctor, error)
	GetDoctorsBySpecialist(context.Context, bson.M) ([]*types.Doctor, error)
	GetDoctorByID(context.Context, string) (*types.Doctor, error)
	GetDoctorByEmail(context.Context, string) (*types.Doctor, error)
	InsertDoctor(context.Context, *types.Doctor) (*types.Doctor, error)
	SetDoctorAppointmentInfo(context.Context, *types.Doctor, string) *types.Doctor
}

type MongoDoctorStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	CheckUpStore
}

func NewMongoDoctorStore(client *mongo.Client) *MongoDoctorStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoDoctorStore{
		client: client,
		coll:   client.Database(dbname).Collection(doctorColl),
	}
}

func (d *MongoDoctorStore) GetDoctors(ctx context.Context) ([]*types.Doctor, error) {

	var doctors []*types.Doctor
	cur, err := d.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &doctors); err != nil {
		return nil, err
	}
	return doctors, nil
}

func (d *MongoDoctorStore) GetDoctorsBySpecialist(ctx context.Context, filter bson.M) ([]*types.Doctor, error) {
	cur, err := d.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var doctors []*types.Doctor
	if err := cur.All(ctx, &doctors); err != nil {
		return nil, err
	}
	return doctors, nil
}

func (d *MongoDoctorStore) GetDoctorByID(ctx context.Context, id string) (*types.Doctor, error) {
	var doctor types.Doctor
	fmt.Println(id)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("invalid doctor id")
		return nil, err
	}
	if err := d.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doctor); err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (d *MongoDoctorStore) InsertDoctor(ctx context.Context, doctor *types.Doctor) (*types.Doctor, error) {
	res, err := d.coll.InsertOne(ctx, doctor)
	if err != nil {
		return nil, err
	}
	doctor.ID = res.InsertedID.(primitive.ObjectID)
	return doctor, nil
}

func (d *MongoDoctorStore) GetDoctorByEmail(ctx context.Context, email string) (*types.Doctor, error) {
	var doctor types.Doctor
	if err := d.coll.FindOne(ctx, bson.M{"email": email}).Decode(&doctor); err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (m *MongoDoctorStore) SetDoctorAppointmentInfo(ctx context.Context, dr *types.Doctor, id string) *types.Doctor {
	var doc types.Doctor
	filter := bson.D{{Key: "_id", Value: dr.ID}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "appointmentInfo.appNumber", Value: dr.AppointmentInfo.AppointmentNumber},
		{Key: "appointmentInfo.appDate", Value: dr.AppointmentInfo.AppointmentDate},
		{Key: "appointmentInfo.patient", Value: id},
	}}}
	if err := m.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		fmt.Println("Error updating doctor appointment info:", err)
		return nil
	}
	return &doc

}
