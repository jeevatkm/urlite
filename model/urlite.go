package model

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//const URLITE_COLLECTION = "multiple since per domain"

type Urlite struct {
	ID         string    `bson:"_id,omitempty"`
	Urlite     string    `bson:"urlite"`
	LongUrl    string    `bson:"long_url"`
	CreateTime time.Time `bson:"ct"`
}

func CreateUrlite(db *mgo.Database, coll string, ul *Urlite) error {
	return db.C(coll).Insert(ul)
}

func GetUrlite(db *mgo.Database, coll string, id *string) (ul *Urlite, err error) {
	err = db.C(coll).Find(bson.M{"_id": id}).One(&ul)
	return
}
