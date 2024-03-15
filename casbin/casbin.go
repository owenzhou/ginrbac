package casbin

import (
	"log"

	"github.com/owenzhou/ginrbac/support/facades"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func newCasbin() *casbin.Enforcer {
	if facades.Config == nil {
		return nil
	}

	if facades.DB == nil {
		log.Println("Casbin error: can not connect to database.")
		return nil
	}

	//使用代码创建模型，不使用conf文件，免得在打包时需要将文件打包进去
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act || r.sub == \"root\"")
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(facades.DB, &casbinModel{})
	if err != nil {
		log.Println("Casbin adapter error: ", err)
		return nil
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Println("Casbin enforcer error: ", err)
		return nil
	}

	return enforcer
}
