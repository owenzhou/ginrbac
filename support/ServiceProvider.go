package support

import (
	"github.com/owenzhou/ginrbac/contracts"
)

//根服务
type ServiceProvider struct {
	App contracts.IApplication
}

func (s *ServiceProvider) Boot() {

}

func (s *ServiceProvider) Register() {

}
