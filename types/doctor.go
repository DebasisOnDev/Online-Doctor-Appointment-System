package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Appointment struct {
	AppointmentNumber int       `json:"appNumber"`
	AppointmentDate   time.Time `json:"appDate"`
	Patient           string    `json:"patient"`
}

type Doctor struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string               `bson:"firstName" json:"firstName"`
	LastName          string               `bson:"lastName" json:"lastName"`
	Email             string               `bson:"email" json:"email"`
	EncryptedPassword string               `bson:"encryptedPassword" json:"-"`
	PassoutYear       string               `bson:"passoutYear" json:"passoutYear"`
	Certificate       string               `bson:"certificate" json:"certificate"`
	Experience        int                  `bson:"experience" json:"experience"`
	PreviousEmployer  string               `bson:"previousEmployer" json:"previousEmployer"`
	Specialist        string               `bson:"specialist" json:"specialist"`
	GovtID            string               `bson:"govtId" json:"govtId"`
	Fee               int                  `bson:"fee" json:"fee"`
	Appointments      []primitive.ObjectID `bson:"appointments" json:"appointments"`
	AppointmentInfo   Appointment          `bson:"appointmentInfo" json:"appointmentInfo"`
	WorkingHour       string               `bson:"workingHour" json:"workingHour"`
}

type CreateDoctorParams struct {
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Certificate      string `json:"certificate"`
	Experience       int    `json:"experience"`
	WorkingHour      string `json:"workingHour"`
	PreviousEmployer string `json:"previousEmployer"`
	Specialist       string `json:"specialist"`
	GovtID           string `json:"govtId"`
}

func (params CreateDoctorParams) ValidateDoctor() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	if len(params.Certificate) < minCertificateLen {
		errors["certificate"] = fmt.Sprintf("certificate length should be at least %d characters", minCertificateLen)
	}
	if params.Experience < minExperienceLen {
		errors["experience"] = fmt.Sprintf("experience must be at least %d characters", minExperienceLen)
	}
	if len(params.PreviousEmployer) < minPrevEmployerLen {
		errors["previousEmployer"] = fmt.Sprintf("previousEmployer must be of %d characters", minPrevEmployerLen)
	}
	if len(params.Specialist) < minSpecialistLen {
		errors["specialist"] = fmt.Sprintf("specialist length should be of %d", minSpecialistLen)
	}
	if len(params.WorkingHour) < minWorkingHourLen {
		errors["workingHour"] = fmt.Sprintf("doctor availability for hour %d must be stated", minWorkingHourLen)
	}
	if len(params.GovtID) < minGovtIdLen {
		errors["govtId"] = fmt.Sprintf("minimun govt ID length should be of %d", minGovtIdLen)
	}
	return errors
}

func NewDoctorFromParams(params CreateDoctorParams) (*Doctor, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	id := primitive.NewObjectID()
	return &Doctor{
		ID:                id,
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
		Certificate:       params.Certificate,
		Experience:        params.Experience,
		PreviousEmployer:  params.PreviousEmployer,
		Specialist:        params.Specialist,
		GovtID:            params.GovtID,
		Fee:               0,
	}, nil
}
