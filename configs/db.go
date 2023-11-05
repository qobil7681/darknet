package configs

import (
	"fmt"
	"os"
)

// DBConfig holds the configuration details for your database connection.
type DBConfig struct {
	Host       string // Database host address or IP.
	Port       string // Database port.
	DbName     string // Name of the database.
	User       string // Database username.
	Password   string // Password for the database user.
	DbURL      string // Full database connection URL
	DbTestName string // Name of the database used for testing.
}

func LoadDatabase() {
	DBHost := os.Getenv("DB_HOST")
	if DBHost == "" {
		DBHost = "localhost"
	}

	DBPort := os.Getenv("DB_PORT")
	if DBPort == "" {
		DBPort = "5432"
	}

	DbName := os.Getenv("DB_NAME")
	if DbName == "" {
		DbName = "postgres"
	}

	DBUser := os.Getenv("DB_USER")
	if DBUser == "" {
		DBUser = "postgres"
	}

	DBPass := os.Getenv("DB_PASS")

	DBTestingName := os.Getenv("DB_TESTING_NAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPass, DBHost, DBPort, DbName)

	cfg.DbConf = DBConfig{
		Host:       DBHost,
		Port:       DBPort,
		DbName:     DbName,
		User:       DBUser,
		Password:   DBPass,
		DbURL:      dbURL,
		DbTestName: DBTestingName,
	}
}
