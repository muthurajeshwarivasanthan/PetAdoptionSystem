package Config

import (
	"fmt"
	"log"

	Models "pet/Models"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB - Initialize DB Connection
func InitDB() {
	viper.AddConfigPath("./Config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	username := viper.GetString("prod.username")
	password := viper.GetString("prod.password")
	dbHost := viper.GetString("prod.db_host")
	dbPort := viper.GetInt("prod.db_port")
	dbName := viper.GetString("prod.db_name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, dbHost, dbPort, dbName)

	var errDB error
	db, errDB = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errDB != nil {
		log.Fatalf("Error connecting to database: %v", errDB)
	}

	// AutoMigrate Seller Model
	err = db.AutoMigrate(&Models.Seller{}, &Models.Buyer{}, &Models.Pet{}, &Models.Adoption{}, &Models.PetHealth{})
	if err != nil {
		log.Fatalf("Error in migration: %v", err)
	}
	log.Println("Database connected successfully!")
}

// GetDB - Return DB instance
func GetDB() *gorm.DB {
	return db
}
