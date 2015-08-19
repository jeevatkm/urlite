package model

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DOMAIN_COLLECTION = "domains"

type Domain struct {
	ID              bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Name            string        `bson:"name" json:"name"`
	Scheme          string        `bson:"scheme" json:"scheme"`
	Salt            string        `bson:"salt" json:"-"`
	LinkCount       int64         `bson:"link_cnt" json:"link_cnt"`
	CustomLinkCount int64         `bson:"custom_link_cnt" json:"custom_link_cnt"`
	IsDefault       bool          `bson:"is_default" json:"is_default"`
	TrackName       string        `bson:"track_name" json:"-"`
	CreatedBy       string        `bson:"cb" json:"-"`
	CreatedTime     time.Time     `bson:"ct" json:"-"`
	UpdatedBy       string        `bson:"ub" json:"-"`
	UpdatedTime     time.Time     `bson:"ut" json:"-"`
}

func CreateDomain(db *mgo.Database, d *Domain) error {
	d.ID = bson.NewObjectId()
	d.CreatedTime = time.Now()
	return db.C(DOMAIN_COLLECTION).Insert(d)
}

func GetDomain(db *mgo.Database, name *string) (d *Domain, err error) {
	err = db.C(DOMAIN_COLLECTION).Find(bson.M{"name": name}).One(&d)
	return
}

func GetDomainById(db *mgo.Database, id bson.ObjectId) (d *Domain, err error) {
	err = db.C(DOMAIN_COLLECTION).FindId(id).One(&d)
	return
}

func GetAllDomain(db *mgo.Database) (domains []Domain, err error) {
	err = db.C(DOMAIN_COLLECTION).Find(bson.M{}).All(&domains)
	return
}

func GetDefaultDomain(db *mgo.Database) (d *Domain, err error) {
	err = db.C(DOMAIN_COLLECTION).Find(bson.M{"is_default": true}).One(&d)
	return
}

func UpdateDomainLinkCount(db *mgo.Database, d *Domain) (err error) {
	sel := bson.M{"_id": d.ID}
	update := bson.M{"$set": bson.M{"link_cnt": d.LinkCount, "custom_link_cnt": d.CustomLinkCount, "ub": "system", "ut": time.Now()}}
	err = db.C(DOMAIN_COLLECTION).Update(sel, update)
	return
}

func (d *Domain) ComposeUrlite(id *string) string {
	return fmt.Sprintf("%v://%v/%v", d.Scheme, d.Name, *id)
}
