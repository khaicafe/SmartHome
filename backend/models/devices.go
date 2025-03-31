package models

import "gorm.io/gorm"

type MappedSwitch struct {
	gorm.Model
	DeviceID string `json:"device_id" gorm:"index"`
	Code     string `json:"code" gorm:"index"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
}

type Setting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex;not null"`
	Value string `gorm:"not null"`
}
