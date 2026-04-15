package models

import (
	"time"

	"github.com/ShreyanshK1103/pipelineTest/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`

}

func ReturnedUser(d database.User) User {

	return User {
		ID: d.ID,
		Email : d.Email,
		CreatedAt: d.CreatedAt,
	}
}