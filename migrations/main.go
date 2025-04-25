package migrations

import (
	"fmt"
	m00 "product_recommendation/migrations/00_add_demo_product"
	"product_recommendation/pkg/logger"
	"product_recommendation/pkg/model"
	"product_recommendation/pkg/types"
	"product_recommendation/pkg/utils/gormigrate"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {

	migrationScripts, migrateOptions := GetMigrationScripts(db)
	if len(migrationScripts) > 0 {
		m := gormigrate.New(db, &migrateOptions, migrationScripts)

		if err := m.Migrate(); err != nil {
			return fmt.Errorf("could not migrate: %v", err)
		}
	}
	migrateLatestModel(db)
	return nil
}

func GetMigrationScripts(db *gorm.DB) ([]*gormigrate.Migration, gormigrate.Options) {
	migrateOptions := gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              255,
		UseTransaction:            false,
		ValidateUnknownMigrations: true,
	}
	var migrationScripts []*gormigrate.Migration
	migrationScripts = append(migrationScripts,
		&m00.Migration,
	)
	return migrationScripts, migrateOptions
}

func migrateLatestModel(db *gorm.DB) {
	migrations := []interface{}{
		&model.Login{},
		&model.Product{},
	}

	for _, m := range migrations {
		err := db.AutoMigrate(m)
		if err != nil {
			logger.Info("migrate error", types.H{
				"err": err,
			})
			fmt.Printf("could not auto migrate: %v", err)
		}
	}
}
