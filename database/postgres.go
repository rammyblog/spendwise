package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {

	user, exist := os.LookupEnv("DB_USER")

	if !exist {
		log.Fatal("DB_USER not set in .env")
		return nil, fmt.Errorf("DB_USER not set in .env")

	}

	pass, exist := os.LookupEnv("DB_PASSWORD")

	if !exist {
		log.Fatal("DB_PASSWORD not set in .env")
		return nil, fmt.Errorf("DB_PASSWORD not set in .env")

	}

	port, exist := os.LookupEnv("DB_PORT")

	if !exist {
		log.Fatal("DB_PASS not set in .env")
		return nil, fmt.Errorf("DB_PORT not set in .env")

	}

	host, exist := os.LookupEnv("DB_HOST")

	if !exist {
		log.Fatal("DB_HOST not set in .env")
		return nil, fmt.Errorf("DB_HOST not set in .env")

	}

	name, exist := os.LookupEnv("DB_NAME")

	if !exist {
		log.Fatal("DB_NAME not set in .env")
		return nil, fmt.Errorf("DB_NAME not set in .env")

	}

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, user, pass, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// MigrateTables(db)
	// if *seed {
	// 	seeder(db)
	// }
	return db, nil

}

// func MigrateTables(db *gorm.DB) {
// 	fmt.Println("Migrating tables")
// 	// Auto migrate tables here
// 	if err := db.AutoMigrate(&models.User{}); err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// 	if err := db.AutoMigrate(&models.Product{}); err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// 	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
// 		log.Fatal(err)
// 		panic(err)
// 	}

// }
