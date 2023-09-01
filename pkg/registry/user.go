package registry

import (
	"github.com/ndodanli/go-clean-architecture/pkg/adapter/controller"
	"github.com/ndodanli/go-clean-architecture/pkg/adapter/repository"
	"github.com/ndodanli/go-clean-architecture/pkg/usecase/usecase"
)

func (r *registry) NewUserController() controller.User {
	u := usecase.NewUserUsecase(
		repository.NewUserRepository(r.db),
		repository.NewDBRepository(r.db),
	)

	return controller.NewUserController(u)
}
