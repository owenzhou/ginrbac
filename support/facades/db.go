package facades

import "gorm.io/gorm"

var DB *gorm.DB

type DBFacade struct {
	*Facade
}

func (db *DBFacade) GetFacadeAccessor() {
	DB = db.App.Make("db").(*gorm.DB)
}
