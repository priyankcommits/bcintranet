package models

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	// User type represents the registered user.
	User struct {
		ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
		UserID      string        `json:"userid"`
		FirstName   string        `json:"firstname"`
		LastName    string        `json:"lastname"`
		Email       string        `json:"email"`
		AccessToken string        `json:"token,omitempty"`
	}
	// Profile type represents the personal data of a user.
	Profile struct {
		ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
		UserID      string        `json:"userid"`
		Age         string        `json:"age"`
		Mobile      string        `json:"mobile"`
		BloodGroup  string        `json:"bloodgroup"`
		Address     string        `json:"address"`
		Description string        `json:"description"`
		GitHub      string        `json:"github"`
	}
)
