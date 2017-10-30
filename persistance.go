package main

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/satori/go.uuid"
)

var db *gorm.DB

func open() {
	var err error
	db, err = gorm.Open("sqlite3", "wish.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Wish{})
	db.AutoMigrate(&User{})
}

func newUserFromDB(user *User) (string, error) {
	if getUserByMailFromDB(user.Mail) != nil {
		return "", errors.New("Le compte existe déjà")
	}
	user.UserID = uuid.NewV4().String()
	db.Create(user)
	return user.UserID, nil
}

func getUserFromDB(userID string) *User {
	user := User{}
	db.First(&user, userID)
	if user.UserID == "" {
		return nil
	}
	return &user
}

func getUserByMailFromDB(mail string) *User {
	user := User{}
	db.Where(&User{Mail: mail}).First(&user)
	if user.UserID == "" {
		return nil
	}
	return &user
}

func getUserByMailAndPasswordFromDB(mail string, password string) *User {
	user := User{}
	db.Where(&User{Mail: mail, Password: password}).First(&user)
	if user.UserID == "" {
		return nil
	}
	return &user
}

func listWishesFromDB(userID string) []Wish {
	whishes := make([]Wish, 0)
	db.Where(&Wish{UserID: userID}).Find(&whishes)
	return whishes
}

func newWishFromDB(userID string, wish *Wish) (*Wish, error) {
	if getUserByMailFromDB(userID) != nil {
		return nil, errors.New("L'utilisateur n'existe pas")
	}
	wish.UserID = userID
	wish.WishID = uuid.NewV4().String()
	db.Create(wish)
	return wish, nil
}

func getWishesByUserAndID(userID string, wishID string) *Wish {
	wish := Wish{}
	db.Where(&Wish{UserID: userID, WishID: wishID}).First(&wish)
	if wish.WishID == "" {
		return nil
	}
	return &wish
}

func updateWishFromDB(userID string, wishID string, wish *Wish) error {
	wishDb := getWishesByUserAndID(userID, wishID)
	if wishDb == nil {
		return errors.New("Le souhait n'existe pas")
	}
	wishDb.Title = wish.Title
	wishDb.Description = wish.Description
	wishDb.Link = wish.Link
	wishDb.Price = wish.Price
	db.Save(wishDb)
	return nil
}

func deleteWishFromDB(userID string, wishID string) error {
	wishDb := getWishesByUserAndID(userID, wishID)
	if wishDb == nil {
		return errors.New("Le souhait n'existe pas")
	}
	db.Delete(wishDb)
	return nil
}
