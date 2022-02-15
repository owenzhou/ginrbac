package database

import (
	"ginrbac/bootstrap/contracts"
	"ginrbac/bootstrap/support"
)

type DBServiceProvider struct {
	*support.ServiceProvider
}

func (db *DBServiceProvider) Register() {
	db.App.Singleton("db", func(app contracts.IApplication) interface{} {
		return newDB()
	})
}
