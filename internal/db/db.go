package db

import (
	"fmt"
	"log"
	"os"

	"go-backend-react-frontend/internal/models"

	cloudsqlconn "cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() *gorm.DB {
	instanceConn := os.Getenv("INSTANCE_CONNECTION_NAME")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbName == "" {
		dbName = "postgres"
	}

	logLevel := logger.Warn
	if os.Getenv("ENV") == "dev" {
		logLevel = logger.Info
	}

	var dialector gorm.Dialector

	if instanceConn == "" {
		// Local: connect to localhost (tunnel handles auth)
		dbPassword := os.Getenv("DB_PASSWORD")
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		var dsn string
		if dbPassword != "" {
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, dbUser, dbPassword, dbName)
		} else {
			dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, dbUser, dbName)
		}
		dialector = postgres.Open(dsn)
	} else {
		// Production: Cloud SQL with IAM auth via Cloud SQL Go Connector
		_, err := pgxv5.RegisterDriver("cloudsql-postgres",
			cloudsqlconn.WithIAMAuthN(),
			cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()),
		)
		if err != nil {
			log.Fatalf("Failed to register Cloud SQL driver: %v", err)
		}
		dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", instanceConn, dbUser, dbName)
		dialector = postgres.New(postgres.Config{
			DriverName: "cloudsql-postgres",
			DSN:        dsn,
		})
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected")
	return DB
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Session{},
		&models.Category{},
		&models.Module{},
		&models.ModulePrerequisite{},
		&models.ModuleTag{},
		&models.Section{},
		&models.Slide{},
		&models.SlideMedia{},
		&models.Quiz{},
		&models.QuizQuestion{},
		&models.QuizOption{},
		&models.QuizAttempt{},
		&models.QuizAnswer{},
		&models.Enrollment{},
		&models.SlideProgress{},
		&models.Certificate{},
		&models.Notification{},
		&models.AuditLog{},
		&models.Comment{},
		&models.PlatformSettings{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated")
}
