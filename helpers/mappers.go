package helpers

import (
	"time"

	"github.com/google/uuid"
	"github.com/mbrunoon/bootdev-chirpy/internal/database"
)

type UserMap struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapUser(dbUser database.User) UserMap {
	user := UserMap{}
	user.ID = dbUser.ID
	user.Email = dbUser.Email.String
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time

	return user
}
