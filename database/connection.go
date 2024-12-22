package database

import (
	"fmt"
	"kgermando/i-pos-restaurant-api/models"
	"kgermando/i-pos-restaurant-api/utils"
	"strconv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
// var DBSQL *sql.DB

func Connect() {
	p := utils.Env("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic("failed to parse database port 😵!")
	}

	DNS := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", utils.Env("DB_HOST"), port, utils.Env("DB_USER"), utils.Env("DB_PASSWORD"), utils.Env("DB_NAME"))
	connection, err := gorm.Open(postgres.Open(DNS), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("Could not connect to the database 😰!")
	}

	DB = connection
	fmt.Println("Database Connected 🎉!") 

	connection.AutoMigrate(
		&models.User{},
		&models.Entreprise{}, 
		&models.Client{},
		&models.Commande{},
		&models.CommandeLine{},
		&models.Contact{},
		&models.Fournisseur{},
		&models.Pos{},
		&models.Product{},
		&models.Stock{},
	
	)

}