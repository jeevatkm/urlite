package model

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jeevatkm/urlite/util"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const USER_COLLECTION = "users"

type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Email       string        `bson:"email"`
	Password    string        `bson:"password"`
	Permissions []string      `bson:"permissions"`
	CreatedBy   bson.ObjectId `bson:"cb"`
	CreatedTime time.Time     `bson:"ct"`
	UpdatedBy   bson.ObjectId `bson:"ub"`
	UpdatedTime time.Time     `bson:"ut"`
}

// func init() {
// 	gob.Register(&User{})
// }

func (u *User) IsPermission(p string) bool {
	return util.Contains(u.Permissions, p)
}

func (u *User) IsAdmin() bool {
	return u.IsPermission("ADMIN")
}

func (u *User) IsBasic() bool {
	return u.IsPermission("BASIC")
}

func AuthenticateUser(db *mgo.Database, email string, password string) (user *User, err error) {
	user, err = GetUserByEmail(db, email)
	if err != nil {
		log.Errorf("Error occurred while retriving user: %v", err)
		return
	}

	result := util.ComparePassword(user.Password, password)
	if !result {
		err = errors.New("Incorrect credentials")
		user = nil
	}
	return
}

func GetUserById(db *mgo.Database, id bson.ObjectId) (user *User, err error) {
	err = db.C(USER_COLLECTION).FindId(id).One(&user)
	return
}

func GetUserByEmail(db *mgo.Database, email string) (user *User, err error) {
	err = db.C(USER_COLLECTION).Find(bson.M{"email": email}).One(&user)
	return
}

func InsertUser(db *mgo.Database, user *User) error {
	user.ID = bson.NewObjectId()
	return db.C(USER_COLLECTION).Insert(user)
}
