package main

import (
	"context"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;primary_key;"`
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserEvent struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uuid.UUID `gorm:"type:uuid"`
}

type UsersRepository struct {
	db *gorm.DB
}

func (repo *UsersRepository) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{}
	result := repo.db.WithContext(ctx).Find(&user, "id = ?", id)
	return user, result.Error
}

func (repo *UsersRepository) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	user := &User{}
	result := repo.db.WithContext(ctx).Delete(&user, "id = ?", id)
	return result.Error
}

func (repo *UsersRepository) SaveUser(ctx context.Context, id uuid.UUID, firstName string, lastName string) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&User{Id: id, FirstName: firstName, LastName: lastName}); result.Error != nil {
			return result.Error
		}
		if result := tx.Create(&UserEvent{UserId: id}); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func (repo *UsersRepository) SelectUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	result := repo.db.WithContext(ctx).Find(&users)
	return users, result.Error
}
