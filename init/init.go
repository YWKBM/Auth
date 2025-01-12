package initData

import (
	"auth/config"
	"auth/repo"
)

type Init struct {
	config config.Config
	repo   repo.Repos
}

func SetInit(config config.Config, repo repo.Repos) *Init {
	return &Init{config: config, repo: repo}
}

func (i *Init) InitData() {
	i.createAdmin(i.config.ADMIN)
}

func (i *Init) createAdmin(config config.AdminConfig) {
	i.repo.Authorization.CreateUser(config.LOGIN, config.PASSWORD, config.EMAIL, "ADMIN")
}
