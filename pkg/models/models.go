package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Email    string    `db:"email" json:"email"`
	Role     string    `db:"user_role" json:"role"`
	Password string    `db:"password_hash" json:"password"`
}

type PVZ struct {
	ID               uuid.UUID `db:"id" json:"id"`
	RegistrationDate time.Time `db:"registration_date" json:"registrationDate"`
	City             string    `db:"city" json:"city"`
}

type Reception struct {
	ID       uuid.UUID `db:"id" json:"id"`
	DateTime time.Time `db:"date_time" json:"dateTime"`
	PvzID    uuid.UUID `db:"pvz_id" json:"pvzID"`
	Status   string    `db:"reception_status" json:"status"`
}

type Product struct {
	ID          uuid.UUID `db:"id" json:"id"`
	DateTime    time.Time `db:"date_time" json:"dateTime"`
	Type        string    `db:"product_type" json:"type"`
	ReceptionID uuid.UUID `db:"reception_id" json:"receptionId"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
