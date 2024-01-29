package database

import (
	"github.com/ChitrangGoyani/go-projects/tree/main/go-jwt-react/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=postgres dbname=jwtproject port=5432 sslmode=disable"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	DB = connection
	connection.AutoMigrate(&models.User{})
}
