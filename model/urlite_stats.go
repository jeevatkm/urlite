package model

import (
	"time"
)

type UrliteStats struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Note        string    `bson:"note" json:"note"`
	Tags        []string  `bson:"tags" json:"tags"`
	Clicks      int       `bson:"clicks" json:"clicks"`
	CreatedBy   string    `bson:"cb" json:"-"`
	CreatedTime time.Time `bson:"ct" json:"-"`
	UpdatedBy   string    `bson:"ub" json:"-"`
	UpdatedTime time.Time `bson:"ut" json:"-"`
}

type UrliteCountryStats struct {
	ID          string    `bson:"_id,omitempty" json:"-"`
	UrliteId    string    `bson:"urlite_id"`
	CountryCode string    `bson:"cc_iso" json:"country_code_iso"`
	Clicks      string    `bson:"clicks" json:"clicks"`
	ClickTime   time.Time `bson:"click_time" json:"click_time"`
}

type UrliteReferrerStats struct {
	ID        string    `bson:"_id,omitempty" json:"-"`
	UrliteId  string    `bson:"urlite_id"`
	URL       string    `bson:"url" json:"url"`
	Host      string    `bson:"host" json:"host"`
	Clicks    string    `bson:"clicks" json:"clicks"`
	ClickTime time.Time `bson:"click_time" json:"click_time"`
}

type UrliteLog struct {
	ID          string    `bson:"_id,omitempty"`
	UrliteId    string    `bson:"urlite_id"`
	Domain      string    `bson:"domain"`
	Refferrer   string    `bson:"referrer"`
	CountryCode string    `bson:"cc_iso"`
	CreatedTime time.Time `bson:"ct"`
}
