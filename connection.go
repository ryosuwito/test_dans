package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func ConnectToDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:Kemong123@tcp(localhost:3306)/test_dans?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	// AutoMigrate the User struct
	db.AutoMigrate(&User{})
	// Create the dummy user
	createDummyUser(db)
	return db, nil
}
func createDummyUser(db *gorm.DB) {
	var count int
	db.Model(&User{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		dummyUser := User{Username: "admin", Password: "password"}
		db.Create(&dummyUser)
	}
}
