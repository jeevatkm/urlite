package model

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const API_USER_COLLECTION = "api_users"

type ApiUser struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	AppName      string        `bson:"app_name"`
	Bearer       string        `bson:"bearer"`
	Domains      []string      `bson:"domains"`
	LastAccessed time.Time     `bson:"last_accessed"`
	IPAddress    string        `bson:"ip_address"`
	CreatedBy    bson.ObjectId `bson:"cb"`
	CreatedTime  time.Time     `bson:"ct"`
	UpdatedBy    bson.ObjectId `bson:"ub"`
	UpdatedTime  time.Time     `bson:"ut"`
}

func CreateApiUser(db *mgo.Database, au *ApiUser) error {
	au.ID = bson.NewObjectId()
	return db.C(API_USER_COLLECTION).Insert(au)
}

func GetApiUserByAppName(db *mgo.Database, appName *string) (au *ApiUser, err error) {
	err = db.C(API_USER_COLLECTION).Find(bson.M{"app_name": appName}).One(&au)
	return
}

func GetApiUserByBearer(db *mgo.Database, bearer *string) (au *ApiUser, err error) {
	err = db.C(API_USER_COLLECTION).Find(bson.M{"bearer": bearer}).One(&au)
	return
}

func GetAllApiUser(db *mgo.Database) (apiUsers []ApiUser, err error) {
	err = db.C(API_USER_COLLECTION).Find(bson.M{"_id": bson.M{"$exist": 1}}).Sort("-last_accessed").All(&apiUsers)
	return
}

func UpdateApiUserLastAccess(db *mgo.Database, au *ApiUser) (err error) {
	sel := bson.M{"_id": au.ID}
	update := bson.M{"$set": bson.M{"last_accessed": au.LastAccessed, "ip_address": au.IPAddress}}
	err = db.C(API_USER_COLLECTION).Update(sel, update)
	return
}
