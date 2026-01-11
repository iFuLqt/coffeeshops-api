package model

type Facility struct {
	ID int64 `gorm:"id"`
	Code string `gorm:"code"`
	Name string `gorm:"name"`
}