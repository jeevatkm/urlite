package model

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Urlite struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Urlite      string    `bson:"urlite" json:"urlite"`
	LongUrl     string    `bson:"long_url" json:"long_url"`
	Domain      string    `bson:"domain" json:"domain,omitempty"`
	CreatedBy   string    `bson:"cb" json:"created_by"`
	CreatedTime time.Time `bson:"ct" json:"created_time"`
}

func CreateUrlite(db *mgo.Database, coll string, ul *Urlite) error {
	return db.C(coll).Insert(ul)
}

func GetUrlite(db *mgo.Database, coll string, id *string) (ul *Urlite, err error) {
	err = db.C(coll).Find(bson.M{"_id": id}).One(&ul)
	return
}

func GetUrliteByPage(db *mgo.Database, coll string, query bson.M, page *Pagination) (*PaginatedResult, error) {
	urlites := []Urlite{}
	total, err := db.C(coll).Find(query).Count()
	err = db.C(coll).Find(query).Sort(page.Sort).Skip(page.Offset).Limit(page.Limit).All(&urlites)

	return &PaginatedResult{Total: total, Result: urlites}, err
}
