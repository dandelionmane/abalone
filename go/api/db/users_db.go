package db

import (
	api "github.com/danmane/abalone/go/api"
	"github.com/jinzhu/gorm"
)

type usersDB struct {
	db *gorm.DB
}

func (s *usersDB) Create(u api.User) (*api.User, error) {
	if err := s.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *usersDB) List() ([]api.User, error) {
	var users []api.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *usersDB) Delete(id int64) error {
	return s.db.Delete(api.User{ID: id}).Error
}

var _ api.UsersService = new(usersDB)