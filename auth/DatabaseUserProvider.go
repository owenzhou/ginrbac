package auth

import (
	"github.com/owenzhou/ginrbac/contracts"

	"gorm.io/gorm"
)

type DatabaseUserProvider struct {
	db     *gorm.DB
	table  string
	hasher contracts.IHash
}

func (d *DatabaseUserProvider) RetrieveById(identifier string) contracts.IUser {
	data := make(map[string]interface{})
	d.db.Table(d.table).Where("id", identifier).Find(&data)
	return &GenericUser{Attributes: data}
}

func (d *DatabaseUserProvider) RetrieveByToken(identifier, token string) contracts.IUser {
	data := make(map[string]interface{})
	d.db.Table(d.table).Where("id", identifier).Find(&data)
	user := &GenericUser{Attributes: data}
	if !user.IsEmpty() && user.GetRememberToken() != "" && user.GetRememberToken() == token {
		return user
	} else {
		return nil
	}
}

func (d *DatabaseUserProvider) UpdateRememberToken(user contracts.IUser, token string) {
	d.db.Table(d.table).
		Where(user.GetAuthIdentifierName(), user.GetAuthIdentifier()).
		Update(user.GetRememberTokenName(), token)
}

func (d *DatabaseUserProvider) UpdateApiToken(user contracts.IUser, token string) {
	d.db.Table(d.table).
		Where(user.GetAuthIdentifierName(), user.GetAuthIdentifier()).
		Update(user.GetApiTokenName(), token)
}

//通过凭证检索
func (d *DatabaseUserProvider) RetrieveByCredentials(credentials map[string]interface{}) contracts.IUser {
	_, ok := credentials["password"]
	if len(credentials) <= 0 || (len(credentials) == 1 && ok) {
		return nil
	}

	condition := make(map[string]interface{})
	for k, v := range credentials {
		if k == "password" {
			continue
		}
		condition[k] = v
	}

	data := make(map[string]interface{})
	d.db.Table(d.table).Where(condition).Take(&data)
	return &GenericUser{Attributes: data}
}

func (d *DatabaseUserProvider) ValidateCredentials(user contracts.IUser, credentials map[string]interface{}) bool {
	return d.hasher.Check(credentials["password"].(string), user.GetAuthPassword())
}
