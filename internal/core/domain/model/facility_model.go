package model

type Facility struct {
	ID int `gorm:"id"`
	Code string `gorm:"code"`
	Name string `gorm:"name"`
}