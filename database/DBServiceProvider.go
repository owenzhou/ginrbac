package database

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support"
)

type DBServiceProvider struct {
	*support.ServiceProvider
}

func (db *DBServiceProvider) Register() {
	db.App.Bind("db", func(app contracts.IApplication) interface{} {
		return newDB()
	})
}
