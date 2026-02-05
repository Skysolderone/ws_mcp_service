package model

import "gorm.io/gorm"

type Rsi struct {
	gorm.Model
	Symbol   string  `gorm:"column:symbol"`
	Value    float64 `gorm:"column:value"`
	Time     string  `gorm:"column:time"` // 格式 2006-01-02 15:04:05
	Interval string  `gorm:"column:interval"`
}

func (Rsi) TableName() string {
	return "rsi"
}
