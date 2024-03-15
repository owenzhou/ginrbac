package auth

import (
	"github.com/owenzhou/ginrbac/contracts"
	"github.com/owenzhou/ginrbac/support/facades"
)

type AuthManager struct {
	App    contracts.IApplication
	guards map[string]contracts.IAuth
}

func newAuthManager(app contracts.IApplication) *AuthManager {
	return &AuthManager{App: app, guards: make(map[string]contracts.IAuth)}
}

// 设置守卫
func (a *AuthManager) Guard(guard ...string) contracts.IAuth {
	var name string
	conf := facades.Config.Auth
	if len(guard) <= 0 {
		name = conf.Defaults["guard"]
	} else {
		name = guard[0]
	}

	if auth, ok := a.guards[name]; ok {
		return auth
	}
	return a.resolve(name)
}

// 守卫的具体实现
func (a *AuthManager) resolve(name string) contracts.IAuth {
	var driver contracts.IAuth
	guards := facades.Config.Auth.Guards

	if _, ok := guards[name]; !ok {
		panic("Can't find guard: " + name)
	}

	switch guards[name]["driver"] {
	case "session":
		driver = a.createSessionDriver(name, guards[name]["provider"])
	case "token":
		driver = a.createTokenDriver(name, guards[name]["provider"])
	case "jwt":
		driver = a.createJWTDriver(name, guards[name]["provider"])
	default:
		panic("Can't find guard: " + name + " driver:" + guards[name]["driver"])
	}
	a.guards[name] = driver
	return driver
}

// 创建session守卫
func (a *AuthManager) createSessionDriver(name, providerName string) contracts.IAuth {

	userProvider := a.createUserProvider(providerName)

	return &SessionGuard{
		name:     name,
		provider: userProvider,
	}
}

// 创建token守卫
func (a *AuthManager) createTokenDriver(name, providerName string) contracts.IAuth {

	userProvider := a.createUserProvider(providerName)

	return &TokenGuard{
		name:     name,
		provider: userProvider,
	}
}

// 创建token守卫
func (a *AuthManager) createJWTDriver(name, providerName string) contracts.IAuth {

	userProvider := a.createUserProvider(providerName)

	return &JWTGuard{
		name:     name,
		provider: userProvider,
	}
}

// 创建用户提供者
func (a *AuthManager) createUserProvider(providerName string) contracts.IUserProvider {
	var userProvider contracts.IUserProvider
	providers := facades.Config.Auth.Providers

	switch providers[providerName]["driver"] {
	case "database":
		userProvider = &DatabaseUserProvider{
			db:     facades.DB,
			table:  providers[providerName]["table"],
			hasher: facades.Hash,
		}
	default:
		panic("Can't find UserProvider: " + providerName + ".")
	}
	return userProvider
}
