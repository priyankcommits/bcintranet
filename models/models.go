package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

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
		UserID     string    `json:"userid"`
		Age        time.Time `json:"age"`
		Mobile     string    `json:"mobile"`
		BloodGroup string    `json:"bloodgroup"`
		Address    string    `json:"address"`
		TagLine    string    `json:"tagline"`
		GitHub     string    `json:"github"`
		SlackName  string    `json:"slackname"`
		PAN        string    `json:"pan"`
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
		Day     time.Time `json:"day"`
		InTime  int       `json:"intime"`
		OutTime int       `json:"outtime"`
		OOO     int       `json:"ooo"`
		WFH     int       `json:"wfh"`
	}
	// Metrics - Daily Logs
	MetricsDailyLogs struct {
		Day         time.Time `json:"day"`
		UserID      string    `json:"userid"`
		Heading     string    `json:"heading"`
		Description string    `json:"text"`
	}
	// Wall Post
	WallPost struct {
		Day      time.Time         `json:"day"`
		PostID   bson.ObjectId     `bson:"_id,omitempty" json:"id"`
		User     User              `json:"user"`
		Type     string            `json:"type"`
		Text     string            `json:"text"`
		Comments []WallPostComment `bson:"comments" json:"comments"`
	}
	// Wall Post - Comment
	WallPostComment struct {
		CommentID bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Day       time.Time     `bson:"day" json:"day"`
		PostID    string        `bson:"post_id" json:"post_id"`
		User      User          `bson:"user" json:"user"`
		Text      string        `bson:"text" json:"text"`
	}
	// Project
	Project struct {
		ProjectID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
		User            User          `bson:"user" json:"user"`
		Title           string        `json:"title"`
		About           string        `json:"about"`
		Aim             string        `json:"aim"`
		TechStack       string        `json:"techstack"`
		CurrentFeatures string        `json:"currentfeatures"`
		FutureFeatures  string        `json:"futurefeatures"`
		Limitations     string        `json:"limitations"`
		RepoLink        string        `json:"repolink"`
		HostLink        string        `json:"hostlink"`
	}
	// About
	About struct {
		AboutID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
		User            User          `bson:"user" json:"user"`
		Title           string        `json:"title"`
		About           string        `json:"about"`
		Aim             string        `json:"aim"`
		TechStack       string        `json:"techstack"`
		CurrentFeatures string        `json:"currentfeatures"`
		FutureFeatures  string        `json:"futurefeatures"`
		Limitations     string        `json:"limitations"`
		RepoLink        string        `json:"repolink"`
		HostLink        string        `json:"hostlink"`
	}
	// Candidate
	Candidate struct {
		CandidateID  bson.ObjectId      `bson:"_id,omitempty" json:"id"`
		Name         string             `json:"name"`
		Email        string             `json:"email"`
		Mobile       string             `json:"mobile"`
		Agency       string             `json:"agency"`
		ProfileLink  string             `json:"link"`
		Status       string             `json:"status"`
		LastUpdateBy User               `bson:"user" json:"user"`
		LastUpdateOn time.Time          `bson:"lastupdate" json:lastupdate`
		Comments     []CandidateComment `bson:"comments" json:"comments"`
	}
	// Candidate Profile Comment
	CandidateComment struct {
		CommentID   bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Day         time.Time     `bson:"day" json:"day"`
		CandidateID string        `bson:"post_id" json:"post_id"`
		User        User          `bson:"user" json:"user"`
		Text        string        `bson:"text" json:"text"`
	}
	// Pay Slip
	Payslip struct {
		PayslipID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
		Name              string        `json:"name"`
		Requestor         User          `bson:"requestor" json:"requestor"`
		Approver          User          `bson:"approver" json:"approver"`
		RequestedOn       time.Time     `json:requestedon`
		Day               time.Time     `json:"day"`
		Month             time.Time     `json:"month"`
		GrossAnnualSalary float64       `json:"salary"`
		TDS               float64       `json:"tds"`
		Amount            float64       `json:"amount"`
		AccountNo         string        `json:"accountno"`
		IFSCCode          string        `json:"ifsccode"`
		Position          string        `json:"position"`
		EmployeeNo        string        `json:"employeeno"`
		Status            int           `json:"status"`
	}
	// Metric ticker notifications
	MetricTicker struct {
		MetricTickerID bson.ObjectId `bson:"_id" json:"id"`
		From           User          `bson:"from" json:"from"`
		To             User          `bson:"to" json:"to"`
		ResourceID     string        `json:"resource"`
		Type           int           `json:"type"`
		Status         int           `json:"status"`
		CreatedOn      time.Time     `json:"createdon"`
	}
	// Notifications
	Notification struct {
		NotificationID bson.ObjectId `bson:"_id" json:"id"`
		From           User          `bson:"from" json:"from"`
		To             User          `bson:"to" json:"to"`
		ResourceID     string        `json:"resource"`
		Type           int           `json:"type"`
		Status         int           `json:"status"`
		CreatedOn      time.Time     `json:"createdon"`
	}
	// Email sendgrid
	Email struct {
		Body    string
		Subject string
		To      User
	}
)
