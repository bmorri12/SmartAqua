// database initial and migrate
package mysql

import (
	"github.com/bmorri12/SmartAqua/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func MigrateDatabase(dbhost, dbport, dbname, dbuser, dbpass string) error {
	mysqldb, err := GetClient(dbhost, dbport, dbname, dbuser, dbpass)
	if err != nil {
		return err
	}
	db, err := gorm.Open("mysql", mysqldb)
	if err != nil {
		return err
	}

	// Then you could invoke `*sql.DB`'s functions with it
	err = db.DB().Ping()
	if err != nil {
		return err
	}

	// Disable table name's pluralization
	db.SingularTable(true)
	db.LogMode(false)

	db.DB().Query("CREATE DATABASE SeaWater; ")
	db.DB().Query("USE SeaWater;")
	// Automating Migration
	db.Set("gorm:table_options", "ENGINE=MyISAM").AutoMigrate(
		&models.Device{},
		&models.Product{},
		&models.Vendor{},
		&models.Application{},
		&models.Rule{},
	)

	return nil
}
