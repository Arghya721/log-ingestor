package config

import (
	"github.com/spf13/viper"
)

// DatabaseConfiguration hold the values required to connect with the database
type DatabaseConfiguration struct {
	DbType      string
	DbUsername  string
	DbPassword  string
	DbName      string
	DbHost      string
	DbPort      string
	DatabaseURL string
	DbSSLMode   string
}

// DatabaseConfig holds the database configurations after reading it from config file
var DatabaseConfig DatabaseConfiguration

func GetDatabaseConfig() {

	// Load the database configuration in the database struct
	DatabaseConfig.DbType = viper.GetString("DB_TYPE")
	DatabaseConfig.DbUsername = viper.GetString("DB_USERNAME")
	DatabaseConfig.DbPassword = viper.GetString("DB_PASSWORD")
	DatabaseConfig.DbName = viper.GetString("DB_NAME")
	DatabaseConfig.DbHost = viper.GetString("DB_HOST")
	DatabaseConfig.DbPort = viper.GetString("DB_PORT")
	DatabaseConfig.DbSSLMode = viper.GetString("DB_SSL_MODE")
}

// GenerateDatabaseURL will generate the url which will be used by our connector
func GenerateDatabaseURL() {

	GetDatabaseConfig()

	if DatabaseConfig.DbType == "postgres" {
		DatabaseConfig.DatabaseURL = "host=" + DatabaseConfig.DbHost + " user=" + DatabaseConfig.DbUsername + " password=" + DatabaseConfig.DbPassword + " dbname=" + DatabaseConfig.DbName + " port=" + DatabaseConfig.DbPort + " sslmode=" + DatabaseConfig.DbSSLMode
	} else if DatabaseConfig.DbType == "mysql" {
		DatabaseConfig.DatabaseURL = DatabaseConfig.DbUsername + ":" + DatabaseConfig.DbPassword + "@tcp(" + DatabaseConfig.DbHost + ":" + DatabaseConfig.DbPort + ")/" + DatabaseConfig.DbName
	} else {
		DatabaseConfig.DatabaseURL = "host=" + DatabaseConfig.DbHost + " user=" + DatabaseConfig.DbUsername + " password=" + DatabaseConfig.DbPassword + " dbname=" + DatabaseConfig.DbName + " port=" + DatabaseConfig.DbPort + " sslmode=" + DatabaseConfig.DbSSLMode
	}
}
