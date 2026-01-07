package entity

type CategoryEntity struct {
	ID int
	Name string
	Slug string
	User UserEntity
}