package store

import (
	"bcintranet/helpers"
	"bcintranet/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetSession(collection string, pk string) *mgo.Session {
	// Dial to database and Return mgo session
	session, err := mgo.Dial("localhost")
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

func SaveUser(userId string, firstName string, lastName string, email string, accessToken string, avatar string) error {
	// Create user data
	session := GetSession("User", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB("bcintranet").C("User")
	_, err := GetUser(userId)
	if err == nil {
		err = c.Update(
			bson.M{"userid": userId},
			bson.M{"$set": bson.M{
				"userid": userId, "firstname": firstName,
				"lastname": lastName, "email": email,
				"accesstoken": accessToken, "avatar": helpers.ImageToBase64(avatar),
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
		err = c.Insert(user)
	}
	return err
}

func GetProfile(userId string) (models.Profile, error) {
	// Find profile
	session := GetSession("Profile", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB("bcintranet").C("Profile")
	var profile models.Profile
	err := c.Find(bson.M{"userid": userId}).One(&profile)
	return profile, err
}

func SaveProfile(profile *models.Profile) error {
	// create user profile
	session := GetSession("Profile", "UserID")
	session = session.Copy()
	defer session.Close()
	c := session.DB("bcintranet").C("Profile")
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
