package users

import (
	"errors"
	"fmt"
	"log"

	"github.com/adamstrickland/dapper-api/internal"
	"github.com/adamstrickland/dapper-api/internal/config"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                  uint   `gorm:"primaryKey"`
	Email               string `gorm:"uniqueIndex"`
	UnencryptedPassword string
	FirstName           string
	LastName            string
}

func Update(cfg *config.Config, u *User) (*User, error) {
	uu, err := FindByEmail(cfg, u.Email)

	if err != nil {
		log.Printf("Unable to find user with email '%s': %e", u.Email, err)
		return nil, err
	}

	uu.FirstName = u.FirstName
	uu.LastName = u.LastName

	db, err := internal.NewConnection(cfg)

	if err != nil {
		log.Printf("Unable to connect to database: %e", err)
		return nil, err
	}

	result := db.Save(uu)

	if result.Error != nil {
		log.Printf("Unable to update User record: %e", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return u, nil
	}

	return uu, nil
}

func Create(cfg *config.Config, u *User) (*User, error) {
	db, err := internal.NewConnection(cfg)

	if err != nil {
		log.Printf("Unable to connect to database: %e", err)
		return nil, err
	}

	query := db.Where(u, "email").Find(&[]User{})

	if query.Error != nil {
		log.Printf("Unable to query User records: %e", query.Error)
		return nil, query.Error
	}

	if query.RowsAffected > 0 {
		return nil, errors.New(fmt.Sprintf("Found extant User record with email '%s'", u.Email))
	}

	result := db.Create(u)

	if result.Error != nil {
		log.Printf("Unable to create User record: %e", result.Error)
		return nil, result.Error
	}

	return u, nil
}

func FindByEmail(cfg *config.Config, email string) (*User, error) {
	db, err := internal.NewConnection(cfg)

	if err != nil {
		log.Printf("Unable to connect to database: %e", err)
		return nil, err
	}

	var user User
	result := db.Where("email = ?", email).Limit(1).Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("No record found")
	}

	return &user, nil
}

func All(cfg *config.Config) (*[]User, error) {
	db, err := internal.NewConnection(cfg)

	if err != nil {
		log.Printf("Unable to connect to database: %e", err)
		return nil, err
	}

	var users []User

	result := db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return &users, nil
}
