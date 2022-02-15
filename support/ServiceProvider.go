package support

import (
	"ginrbac/bootstrap/contracts"
)

//根服务
type ServiceProvider struct {
	App contracts.IApplication
}

func (s *ServiceProvider) Boot() {

}

func (s *ServiceProvider) Register() {

}
