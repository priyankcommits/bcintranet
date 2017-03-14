package models

type (
	// User type represents the registered user.
	User struct {
		UserID      string `json:"userid"`
		FirstName   string `json:"firstname"`
		LastName    string `json:"lastname"`
		Email       string `json:"email"`
		AccessToken string `json:"token,omitempty"`
		Avatar      string `json:"avatar"`
		Admin       bool   `json:"useradmin"`
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
		SlackName  string `json:"slackname"`
	}
	// Flash message Struct
	Message struct {
		Value string
	}
	// Pass keyword args to urls
	Kwargs struct {
		Key   string
		Value string
	}
	// Metrics - Storing Attendance Metrics by day
	MetricsAttendance struct {
		Day   string `json:"day"`
		Month int    `json:"month"`
		Value int    `json:"value"`
	}
	// Metrics - Daily Logs
	MetricsDailyLogs struct {
		Day     string `json:"day"`
		UserID  string `json:"userid"`
		Heading string `json:"heading"`
		Text    string `json:"text"`
	}
)
