package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection using environment variables.
// Ensure LoadEnv() is called before this function (if using a .env file locally),
// and that the necessary environment variables (DB_HOST, DB_USER, DB_PASSWORD,
// DB_NAME, DB_PORT) are set. In a Docker Compose setup, these are typically
// provided by the docker-compose.yml file.
func InitDB() (*gorm.DB, error) {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbTimeZone := os.Getenv("DB_TIMEZONE")

	// Validate essential environment variables
	if dbHost == "" {
		log.Fatal("DB_HOST environment variable is not set. In Docker Compose, this is often the service name of your PostgreSQL container (e.g., 'postgres' or 'db').")
	}
	if dbUser == "" {
		log.Fatal("DB_USER environment variable is not set.")
	}
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set.")
	}
	if dbPort == "" {
		dbPort = "5432" // Default PostgreSQL port
		log.Println("DB_PORT environment variable not set, defaulting to 5432.")
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable" // Common default for development, consider 'require' or 'verify-full' for production
	}
	if dbTimeZone == "" {
		dbTimeZone = "UTC" // Good default
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone)

	log.Printf("Attempting to connect to PostgreSQL with DSN: host=%s user=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbName, dbPort, dbSSLMode, dbTimeZone) // Avoid logging password directly in production logs

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// Avoid printing the full DSN with password in the panic message for security
		return nil, fmt.Errorf("failed to connect to PostgreSQL database. Host: %s, DBName: %s. Error: %v", dbHost, dbName, err)
	}
	log.Println("Successfully connected to PostgreSQL database!")
	return db, nil
}
