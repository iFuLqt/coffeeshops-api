package entity

type CategoryEntity struct {
	ID int64
	Name string
	Slug string
	CreatedBy UserEntity
}