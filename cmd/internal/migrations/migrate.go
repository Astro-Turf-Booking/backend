package migrations

import (
	"log"

	"turfzone-backend/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		// Register future models here:
		// &models.Organization{},
		// &models.OrganizationMember{},
		// &models.Turf{},
		// &models.Booking{},
		// &models.Payment{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
