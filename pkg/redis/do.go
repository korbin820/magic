package redis

import (
	util "service-tag/utils"
	"time"
)

type Redistag struct {
	ID         uint           `gorm:"column:id;primaryKey" json:"id"`
	Module     string         `gorm:"column:module;not null" json:"module"`
	Project    string         `gorm:"column:project;not null" json:"project"`
	Instance   string         `gorm:"column:instance" json:"instance"`
	Updatetime *util.DateTime `gorm:"column:updatetime" json:"updatetime"`
}

type Redis struct {
	ID         uint      `gorm:"column:id"`
	Url        string    `gorm:"column:url;not null"`
	Name       string    `gorm:"column:name;not null"`
	Createtime time.Time `gorm:"autoCreateTime;column:createtime;type:date"`
}
