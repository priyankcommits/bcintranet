package store

import (
	"bcintranet/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetSession(collection string, pk string) *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(collection, pk, session)
	return session
}

func ensureIndex(collection string, pk string, s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("bcintranet").C(collection)
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
	c := session.DB("bcintranet").C("User")
	var user models.User
	err := c.Find(bson.M{"userid": userId}).One(&user)
	return user, err
}

func CreateUserData(userId string, firstName string, lastName string, email string, accessToken string) error {
	// Create user data on auth callback
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB("bcintranet").C("User")
	var user models.User
	user.UserID = userId
	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.AccessToken = accessToken
	err := c.Insert(user)
	return err
}

func GetProfile(userId string) error {
	// Find profile
	session := GetSession("Profile", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB("bcintranet").C("Profile")
	var profile models.Profile
	err := c.Find(bson.M{"userid": userId}).One(&profile)
	return err
}

func CreateProfile(profile *models.Profile) error {
	session := GetSession("Profile", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB("bcintranet").C("Profile")
	err := c.Insert(&profile)
	return err
}
