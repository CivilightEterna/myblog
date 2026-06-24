package models

import "time"

type SiteSetting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SettingKey   string    `gorm:"uniqueIndex;not null" json:"setting_key"`
	SettingValue string    `gorm:"not null" json:"setting_value"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SiteSetting) TableName() string {
	return "site_setting"
}
