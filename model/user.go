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
	ID              bson.ObjectId `bson:"_id,omitempty"`
	Email           string        `bson:"email"`
	Password        string        `bson:"password"`
	IsActive        bool          `bson:"is_active"`
	Permissions     []string      `bson:"permissions"`
	Bearer          string        `bson:"bearer"`
	Domains         []string      `bson:"domains"`
	LastLoggedIn    time.Time     `bson:"last_logged_in"`
	LoginIPAddress  string        `bson:"login_ip_address"`
	LastApiAccessed time.Time     `bson:"last_api_accessed"`
	ApiIPAddress    string        `bson:"api_ip_address"`
	CreatedBy       string        `bson:"cb"`
	CreatedTime     time.Time     `bson:"ct"`
	UpdatedBy       string        `bson:"ub"`
	UpdatedTime     time.Time     `bson:"ut"`
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

func GetActiveUserCount(db *mgo.Database) int {
	cnt, err := db.C(USER_COLLECTION).Find(bson.M{"is_active": true}).Count()
	if err != nil {
		log.Errorf("Unable to get active user count: %q", err)
		return 0
	}
	return cnt
}

func GetUserCount(db *mgo.Database) int {
	cnt, err := db.C(USER_COLLECTION).Count()
	if err != nil {
		log.Errorf("Unable to get user count: %q", err)
		return 0
	}
	return cnt
}

func GetUserById(db *mgo.Database, id bson.ObjectId) (user *User, err error) {
	err = db.C(USER_COLLECTION).FindId(id).One(&user)
	return
}

func GetUserByEmail(db *mgo.Database, email string) (user *User, err error) {
	err = db.C(USER_COLLECTION).Find(bson.M{"email": email}).One(&user)
	return
}

func GetUserByBearer(db *mgo.Database, bearer *string) (user *User, err error) {
	err = db.C(USER_COLLECTION).Find(bson.M{"bearer": bearer}).One(&user)
	return
}

func GetAllUsers(db *mgo.Database) (users []User, err error) {
	err = db.C(USER_COLLECTION).Find(bson.M{}).Sort("-last_api_accessed", "-last_logged_in").All(&users)
	return
}

func CreateUser(db *mgo.Database, user *User) error {
	user.ID = bson.NewObjectId()
	return db.C(USER_COLLECTION).Insert(user)
}

func UpdateUserLastApiAccess(db *mgo.Database, u *User) (err error) {
	sel := bson.M{"_id": u.ID}
	update := bson.M{"$set": bson.M{"last_api_accessed": u.LastApiAccessed, "api_ip_address": u.ApiIPAddress}}
	err = db.C(USER_COLLECTION).Update(sel, update)
	return
}

func UpdateUserLastLogin(db *mgo.Database, u *User) (err error) {
	sel := bson.M{"_id": u.ID}
	update := bson.M{"$set": bson.M{"last_logged_in": u.LastLoggedIn, "login_ip_address": u.LoginIPAddress}}
	err = db.C(USER_COLLECTION).Update(sel, update)
	return
}
