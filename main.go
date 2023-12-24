package main

import (
	"context"
	"log"
	"os"

	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/api"
	"github.com/DebasisOnDev/Online-Doctor-Appointment-System/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoendpoint := os.Getenv("MONGODB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoendpoint))
	if err != nil {
		log.Fatal(err)
	}

	//Handlers initialization
	var (
		userStore    = db.NewMongoUserStore(client)
		doctorStore  = db.NewMongoDoctorStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		checkupStore = db.NewMongoCheckUpStore(client)
		store        = &db.Store{
			User:    userStore,
			Doctor:  doctorStore,
			Booking: bookingStore,
			CheckUp: checkupStore,
		}
		bookingHandler = api.NewBookingHandler(store)
		userHandler    = api.NewUserHandler(userStore)
		doctorHandler  = api.NewDoctorHandler(doctorStore, checkupStore)
		checkupHandler = api.NewCheckUpHandler(store)
		authHandler    = api.NewAuthHandler(userStore, doctorStore)
		app            = fiber.New(config)
		auth           = app.Group("/api/auth")
		dapiv1         = app.Group("/api/dv1", api.JWTAuthenticationDoctor(doctorStore))
		uapiv1         = app.Group("/api/uv1", api.JWTAuthenticationUser(userStore))
	)

	//versioned api routes and auth routes
	//initial phase routes should be replaced and handled with advanced routes in later time
	auth.Post("/user/register", authHandler.HandleRegisterUser)
	auth.Post("/user/login", authHandler.HandleLoginUser)
	auth.Post("/user/logout", authHandler.HandleUserLogOut)

	auth.Post("/doctor/register", authHandler.HandleRegisterDoctor)
	auth.Post("/doctor/login", authHandler.HandleLoginDoctor)
	auth.Post("/doctor/logout", authHandler.HandleDoctorLogOut)

	uapiv1.Get("/user/:id", userHandler.HandleGetUser)
	uapiv1.Get("/user", userHandler.HandleGetUsers)
	uapiv1.Get("/user/doctor", bookingHandler.HandleGetDoctorsByUser)
	//uapiv1.Get("/user/doctor/:specialist", bookingHandler.HandleGetDoctorBySpecialist)
	uapiv1.Get("/user/doctor/:id", bookingHandler.HandleGetDoctorById)
	uapiv1.Get("/user/doctor/:id/book", bookingHandler.HandleBookingADoctorByUser)

	uapiv1.Post("/user", userHandler.HandleInsertUser) //->should be replaced by register user

	dapiv1.Get("/doctor", doctorHandler.HandleGetDoctors)
	dapiv1.Get("/doctor/:id", doctorHandler.HandleGetDoctor)
	dapiv1.Get("/doctor/checkup", doctorHandler.GetCheckups)
	dapiv1.Get("/doctor/checkup/go", checkupHandler.HandleDoCheckUp)
	dapiv1.Get("/doctor/checkup/:id", checkupHandler.HandleGetCheckUpByID)

	//apiv1.Post("/doctor", doctorHandler.HandleInsertDoctor) //->replaced by doctor register

	listenaddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenaddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
