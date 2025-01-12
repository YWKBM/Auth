package initData

import (
	"auth/config"
	"auth/repo"
	"auth/utils"
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
	i.repo.Authorization.CreateUser(config.LOGIN, utils.GnerateHashPassword(config.PASSWORD), config.EMAIL, "ADMIN")
}
