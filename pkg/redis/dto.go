package redis

import (
	time2 "magic/basic/time"
	"time"
)

type Redistag struct {
	ID         uint            `gorm:"column:id;primaryKey" json:"id"`
	Module     string          `gorm:"column:module;not null" json:"module"`
	Project    string          `gorm:"column:project;not null" json:"project"`
	Instance   string          `gorm:"column:instance" json:"instance"`
	Updatetime *time2.DateTime `gorm:"column:updatetime" json:"updatetime"`
}

type Redis struct {
	ID         uint      `gorm:"column:id"`
	Url        string    `gorm:"column:url;not null"`
	Name       string    `gorm:"column:name;not null"`
	Createtime time.Time `gorm:"autoCreateTime;column:createtime;type:date"`
}
