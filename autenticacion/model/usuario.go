package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson" // 👈 solo este import
)

type Usuario struct {
	User     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsuarioDb struct {
	ID           bson.ObjectID `json:"id,omitempty"     bson:"_id,omitempty"` // 👈 bson.ObjectID en vez de primitive.ObjectID
	Username     string        `json:"username"          bson:"username"`
	Email        string        `json:"email"             bson:"email"`
	PasswordHash string        `json:"passwordHash"      bson:"passwordHash"`
	Role         string        `json:"role"              bson:"role"`
	IsActive     bool          `json:"isActive"          bson:"isActive"`
	CreatedAt    time.Time     `json:"createdAt"         bson:"createdAt"`
}
