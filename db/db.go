package db

type Store struct {
	User    UserStore
	Doctor  DoctorStore
	Booking BookingStore
	CheckUp CheckUpStore
}

const MongoDBNameEnvName = "MONGODB_NAME"
