package models

import ()

type (
	// User type represents the registered user.
	User struct {
		UserID      string `json:"userid"`
		FirstName   string `json:"firstname"`
		LastName    string `json:"lastname"`
		Email       string `json:"email"`
		AccessToken string `json:"token,omitempty"`
	}
	// Profile type represents the personal data of a user.
	Profile struct {
		UserID     string `json:"userid"`
		Age        string `json:"age"`
		Mobile     string `json:"mobile"`
		BloodGroup string `json:"bloodgroup"`
		Address    string `json:"address"`
		TagLine    string `json:"tagline"`
		GitHub     string `json:"github"`
	}
	// Flash message Struct
	Message struct {
		Value string
	}
)
