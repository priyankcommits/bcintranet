package store

import (
	"os"
	"time"

	"bcintranet/helpers"
	"bcintranet/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetSession(collection string, pk string) *mgo.Session {
	// Dial to database and Return mgo session
	var session *mgo.Session
	var err error
	if os.Getenv("bc_env") == "development" {
		session, err = mgo.Dial("localhost")
	} else {
		session, err = mgo.Dial(os.Getenv("MONGODB_URI"))
	}
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(collection, pk, session)
	return session
}

func ensureIndex(collection string, pk string, s *mgo.Session) {
	// Ensure an index on the collection, why?
	session := s.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C(collection)
	index := mgo.Index{
		Key:        []string{pk},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func GetUser(userId string) (models.User, error) {
	// get user data
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("User")
	var user models.User
	err := c.Find(bson.M{"userid": userId}).One(&user)
	return user, err
}

func GetAllUsers() ([]models.User, error) {
	// get all users
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("User")
	var users []models.User
	err := c.Find(bson.M{}).All(&users)
	return users, err
}

func GetAdmins() ([]models.User, error) {
	// get all admins only
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("User")
	var users []models.User
	err := c.Find(bson.M{"admin": true}).All(&users)
	return users, err
}

func IsAdmin(userId string) bool {
	// Return admin true or false
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("User")
	var user models.User
	_ = c.Find(bson.M{"userid": userId}).One(&user)
	return user.Admin
}

func SaveUser(userId string, firstName string, lastName string, email string, accessToken string, avatar string, userAdmin bool) error {
	// Create user data
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("User")
	_, err := GetUser(userId)
	if err == nil {
		err = c.Update(
			bson.M{"userid": userId},
			bson.M{"$set": bson.M{
				"userid": userId, "firstname": firstName,
				"lastname": lastName, "email": email,
				"accesstoken": accessToken, "avatar": helpers.ImageToBase64(avatar),
				"useradmin": userAdmin,
			}},
		)
	} else {
		var user models.User
		user.UserID = userId
		user.FirstName = firstName
		user.LastName = lastName
		user.Email = email
		user.AccessToken = accessToken
		user.Avatar = helpers.ImageToBase64(avatar)
		user.Admin = userAdmin
		err = c.Insert(user)
	}
	return err
}

func GetProfile(userId string) (models.Profile, error) {
	// Find profile
	session := GetSession("Profile", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Profile")
	var profile models.Profile
	err := c.Find(bson.M{"userid": userId}).One(&profile)
	return profile, err
}

func SaveProfile(profile *models.Profile) error {
	// create user profile
	session := GetSession("Profile", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Profile")
	_, err := GetProfile(profile.UserID)
	if err == nil {
		err = c.Update(
			bson.M{"userid": profile.UserID},
			bson.M{"$set": &profile},
		)
	} else {
		err = c.Insert(&profile)
	}
	return err
}

func SaveAttendanceLog(attendanceLog *models.MetricsAttendance) error {
	// Save Metric Entry
	session := GetSession("MetricsAttendance", "Day")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("MetricsAttendance")
	var attendanceLogDB models.MetricsAttendance
	err := c.Find(bson.M{"day": attendanceLog.Day}).One(&attendanceLogDB)
	if err == nil {
		err = c.Update(
			bson.M{"day": attendanceLog.Day},
			bson.M{"$set": bson.M{
				"day":    attendanceLog.Day,
				"intime": attendanceLog.InTime, "outtime": attendanceLog.OutTime,
			}},
		)
	} else {
		err = c.Insert(&attendanceLog)
	}
	return err
}

func GetAttendanceMetrics() (models.MetricsAttendance, []models.MetricsAttendance, []models.MetricsAttendance, error) {
	// Get Attendance Metric for pie charts
	session := GetSession("MetricsAttendance", "Day")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("MetricsAttendance")
	var attendanceDayLogDB models.MetricsAttendance
	currentYear, currentMonth, currenDay := time.Now().Date()
	todayDate := time.Date(currentYear, currentMonth, currenDay, 0, 0, 0, 0, time.UTC)
	err := c.Find(bson.M{"day": todayDate}).One(&attendanceDayLogDB)
	var attendanceMonthLogDB []models.MetricsAttendance
	monthStartDate := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	monthEndDate := time.Date(currentYear, currentMonth, 31, 0, 0, 0, 0, time.UTC)
	err = c.Find(bson.M{"day": bson.M{"$gt": monthStartDate, "$lt": monthEndDate}}).All(&attendanceMonthLogDB)
	var attendanceYearLogDB []models.MetricsAttendance
	yearStartDate := time.Date(currentYear, 0, 1, 0, 0, 0, 0, time.UTC)
	yearEndDate := time.Date(currentYear, 12, 31, 0, 0, 0, 0, time.UTC)
	err = c.Find(bson.M{"day": bson.M{"$gt": yearStartDate, "$lt": yearEndDate}}).All(&attendanceYearLogDB)
	return attendanceDayLogDB, attendanceMonthLogDB, attendanceYearLogDB, err
}

func SaveDailyLog(dailyLog *models.MetricsDailyLogs) error {
	// Save Daily Log
	session := GetSession("MetricsDailyLogs", "Day")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("MetricsDailyLogs")
	err := c.Insert(&dailyLog)
	return err

}

func GetDailyLogMetrics() (int, int, int, error) {
	// Get Daily Log Metrics
	session := GetSession("MetricsDailyLogs", "Day")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("MetricsDailyLogs")
	var dayLogCount int
	currentYear, currentMonth, currenDay := time.Now().Date()
	todayDate := time.Date(currentYear, currentMonth, currenDay, 0, 0, 0, 0, time.UTC)
	dayLogCount, err := c.Find(bson.M{"day": todayDate}).Count()
	var monthLogCount int
	monthStartDate := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	monthEndDate := time.Date(currentYear, currentMonth, 31, 0, 0, 0, 0, time.UTC)
	monthLogCount, err = c.Find(bson.M{"day": bson.M{"$gt": monthStartDate, "$lt": monthEndDate}}).Count()
	var yearLogCount int
	yearStartDate := time.Date(currentYear, 0, 1, 0, 0, 0, 0, time.UTC)
	yearEndDate := time.Date(currentYear, 12, 31, 0, 0, 0, 0, time.UTC)
	yearLogCount, err = c.Find(bson.M{"day": bson.M{"$gt": yearStartDate, "$lt": yearEndDate}}).Count()
	return dayLogCount, monthLogCount, yearLogCount, err
}

func GetUserDailyLogs(userId string) ([]models.MetricsDailyLogs, error) {
	// Get all user logs
	session := GetSession("MetricsDailyLogs", "Day")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("MetricsDailyLogs")
	var logs []models.MetricsDailyLogs
	err := c.Find(bson.M{"userid": userId}).Sort("-day").All(&logs)
	return logs, err
}

func SavePost(post *models.WallPost) (*models.WallPost, error) {
	// Save Post
	session := GetSession("WallPost", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("WallPost")
	post.PostID = bson.NewObjectId()
	err := c.Insert(&post)
	return post, err
}

func GetAllPosts() ([]models.WallPost, error) {
	// Get all wall posts
	session := GetSession("WallPost", "User")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("WallPost")
	var posts []models.WallPost
	err := c.Find(nil).Sort("-day").All(&posts)
	return posts, err
}

func GetPost(postId string) (models.WallPost, error) {
	//Get a wall Post
	session := GetSession("WallPost", "User")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("WallPost")
	var posts models.WallPost
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(postId)}).One(&posts)
	return posts, err
}

func DeletePost(postId string) error {
	// Delete a wall post
	session := GetSession("WallPost", "User")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("WallPost")
	err := c.Remove(bson.M{"_id": bson.ObjectIdHex(postId)})
	return err
}

func GetAnnouncements() ([]models.WallPost, error) {
	// Get only announcement posts
	session := GetSession("WallPost", "User")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("WallPost")
	var announcements []models.WallPost
	err := c.Find(bson.M{"type": "2"}).All(&announcements)
	return announcements, err
}

func SaveComment(comment *models.WallPostComment) error {
	// Save Comment
	session := GetSession("WallPost", "User")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("WallPost")
	comment.CommentID = bson.NewObjectId()
	newComment := bson.M{"$push": bson.M{"comments": &comment}}
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(comment.PostID)}, newComment)
	return err
}

func DeleteComment(postId string, commentId string) error {
	// Delete a Comment
	session := GetSession("WallPost", "User")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("WallPost")
	deleteComment := bson.M{"$pull": bson.M{"comments": bson.M{"_id": bson.ObjectIdHex(commentId)}}}
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(postId)}, deleteComment)
	return err
}

func SaveProject(project *models.Project) (*models.Project, error) {
	// Add a project
	session := GetSession("Project", "ProjectID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Project")
	project.ProjectID = bson.NewObjectId()
	err := c.Insert(&project)
	return project, err
}

func GetProject(projectId string) (models.Project, error) {
	// Get a project
	session := GetSession("Project", "ProjectID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Project")
	var project models.Project
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(projectId)}).One(&project)
	return project, err
}

func GetProjects() ([]models.Project, error) {
	// get all projects
	session := GetSession("Project", "ProjectID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Project")
	var projects []models.Project
	err := c.Find(bson.M{}).All(&projects)
	return projects, err
}

func UpdateProject(projectId string, project models.Project) error {
	// Update a project
	session := GetSession("Project", "ProjectID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Project")
	err := c.Update(
		bson.M{"_id": bson.ObjectIdHex(projectId)},
		bson.M{"$set": &project},
	)
	return err
}

func SaveAbout(about *models.About) error {
	// Save about project info
	session := GetSession("About", "AboutID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("About")
	err := c.Insert(&about)
	return err
}

func GetAbout() (models.About, error) {
	// Get about project info
	session := GetSession("About", "AboutID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("About")
	var about models.About
	err := c.Find(bson.M{}).One(&about)
	return about, err
}

func AddCandidate(candidate *models.Candidate) (*models.Candidate, error) {
	// Add a candidate
	session := GetSession("Candidate", "CandidateID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Candidate")
	candidate.CandidateID = bson.NewObjectId()
	err := c.Insert(&candidate)
	return candidate, err
}

func UpdateCandidate(candidateId string, candidate models.Candidate) error {
	// Update a candidate
	session := GetSession("Candidate", "CandidateID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Candidate")
	err := c.Update(
		bson.M{"_id": bson.ObjectIdHex(candidateId)},
		bson.M{"$set": &candidate},
	)
	return err
}

func GetCandidate(candidateId string) (models.Candidate, error) {
	//Get a candidate
	session := GetSession("Candidate", "CandidateID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Candidate")
	var candidate models.Candidate
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(candidateId)}).One(&candidate)
	return candidate, err
}

func GetCandidates() ([]models.Candidate, error) {
	// Get all candidates
	session := GetSession("Candidate", "CandidateID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Candidate")
	var candidates []models.Candidate
	err := c.Find(bson.M{}).All(&candidates)
	return candidates, err
}

func CheckCandidate(email string) error {
	// check if candidate exists
	session := GetSession("Candidate", "CandidateID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Candidate")
	var candidate models.Candidate
	err := c.Find(bson.M{"email": email}).One(&candidate)
	return err
}

func SaveCandidateComment(comment *models.CandidateComment) error {
	// Save Candidate Comment
	session := GetSession("Candidate", "CandidateID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Candidate")
	comment.CommentID = bson.NewObjectId()
	newComment := bson.M{"$push": bson.M{"comments": &comment}}
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(comment.CandidateID)}, newComment)
	return err
}

func SavePayslipRequest(payslip *models.Payslip) error {
	// Save Payslip Request
	session := GetSession("Payslip", "PayslipID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Payslip")
	payslip.PayslipID = bson.NewObjectId()
	err := c.Insert(&payslip)
	return err
}

func GetPayslipHistory(userId string) ([]models.Payslip, error) {
	// Get all payslip requests/approvals for user
	session := GetSession("Payslip", "PayslipID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Payslip")
	var payslips []models.Payslip
	err := c.Find(bson.M{"requestor.userid": userId}).All(&payslips)
	return payslips, err
}

func GetAllPayslips() ([]models.Payslip, error) {
	// Get all payslips
	session := GetSession("Payslip", "PayslipID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Payslip")
	var payslips []models.Payslip
	err := c.Find(bson.M{}).All(&payslips)
	return payslips, err
}

func GetPayslip(payslipId string) (models.Payslip, error) {
	// Get a payslip
	session := GetSession("Payslip", "PayslipID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Payslip")
	var payslip models.Payslip
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(payslipId)}).One(&payslip)
	return payslip, err
}

func UpdatePayslip(payslipId string, payslip models.Payslip) error {
	// Update a payslip
	session := GetSession("Payslip", "PayslipID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Payslip")
	err := c.Update(
		bson.M{"_id": bson.ObjectIdHex(payslipId)},
		bson.M{"$set": &payslip},
	)
	return err
}

func GetNotifications(userId string) ([]models.Notification, error) {
	// Get all notifications for a user
	session := GetSession("Notification", "NotificationID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Notification")
	var notifications []models.Notification
	err := c.Find(bson.M{"to.userid": userId}).Sort("-createdon").All(&notifications)
	return notifications, err
}

func GetMetricTicker(userId string) ([]models.MetricTicker, error) {
	// Get all metric ticker notifications for a user
	session := GetSession("MetricTicker", "MetricTickerID")
	session = session.Copy()
	defer session.Close()
	c := session.DB(os.Getenv("bc_mongo_db")).C("MetricTicker")
	var tickers []models.MetricTicker
	err := c.Find(bson.M{"to.userid": userId}).Sort("-createdon").All(&tickers)
	return tickers, err
}

func CreateNotification(notification *models.Notification) error {
	// Create a notification
	session := GetSession("Notification", "NotificationID")
	session = session.Copy()
	defer session.Close()
	notification.NotificationID = bson.NewObjectId()
	c := session.DB(os.Getenv("bc_mongo_db")).C("Notification")
	err := c.Insert(&notification)
	return err
}

func CreateMetricTicker(ticker *models.MetricTicker) error {
	// Create a metric ticker
	session := GetSession("MetricTicker", "MetricTicker")
	session = session.Copy()
	defer session.Close()
	ticker.MetricTickerID = bson.NewObjectId()
	c := session.DB(os.Getenv("bc_mongo_db")).C("MetricTicker")
	err := c.Insert(&ticker)
	return err
}
