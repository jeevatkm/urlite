package model

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DOMAIN_COLLECTION = "domains"

type Domain struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`
	Scheme     string        `bson:"scheme"`
	Salt       string        `bson:"salt"`
	Count      int64         `bson:"count"`
	CreateTime time.Time     `bson:"ct"`
	UpdateTime time.Time     `bson:"ut"`
}

func CreateDomain(db *mgo.Database, d *Domain) error {
	d.ID = bson.NewObjectId()
	return db.C(DOMAIN_COLLECTION).Insert(d)
}

func GetDomain(db *mgo.Database, name *string) (d *Domain, err error) {
	err = db.C(DOMAIN_COLLECTION).Find(bson.M{"name": name}).One(&d)
	return
}

func GetAllDomain(db *mgo.Database) (domains []Domain, err error) {
	err = db.C(DOMAIN_COLLECTION).Find(bson.M{"_id": bson.M{"$exist": 1}}).All(&domains)
	return
}
