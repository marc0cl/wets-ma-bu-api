package config

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds application configuration
type Config struct {
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDB       string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION", "24"))

	// Use explicit MySQL connection details from environment variables
	host := getEnv("MYSQL_HOST", "")
	port := getEnv("MYSQL_PORT", "")
	user := getEnv("MYSQL_USER", "")
	password := getEnv("MYSQL_PASSWORD", "")
	dbname := getEnv("MYSQL_DB", "")

	return &Config{
		MySQLHost:     host,
		MySQLPort:     port,
		MySQLUser:     user,
		MySQLPassword: password,
		MySQLDB:       dbname,
		DatabaseURL:   "",
		JWTSecret:     getEnv("JWT_SECRET", "your_secret_key"),
		JWTExpiration: jwtExpiration,
	}
}

// InitDatabase initializes and returns a database connection
func (c *Config) InitDatabase() (*gorm.DB, error) {
	// Format MySQL DSN: username:password@tcp(host:port)/dbname?parseTime=true&tls=false
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=false",
		c.MySQLUser, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDB)

	fmt.Printf("Connecting to MySQL database: %s:%s/%s\n", c.MySQLHost, c.MySQLPort, c.MySQLDB)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
