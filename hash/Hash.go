package hash

import (
	"github.com/owenzhou/ginrbac/contracts"

	"golang.org/x/crypto/bcrypt"
)

func newHash(app contracts.IApplication) Hash {
	return Hash{App: app, cost: 10}
}

type Hash struct {
	App  contracts.IApplication
	cost int
}

func (h Hash) Make(value string, options ...map[string]int) string {
	var cost int
	if len(options) <= 0 {
		cost = h.cost
	} else {
		if v, ok := options[0]["rounds"]; !ok {
			cost = h.cost
		} else {
			cost = v
		}
	}
	hashed, err := password_hash(value, cost)
	if err != nil {
		panic("Hash err")
	}
	return hashed
}

func (h Hash) Check(value, hashedValue string, options ...string) bool {
	return password_verify(value, hashedValue)
}

func password_hash(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func password_verify(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
