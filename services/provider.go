package services

import (
	"auth/queue"
	"auth/queue/messages"
	"auth/queue/messages/dto"
	"auth/repo"
	"auth/utils"
	"encoding/json"
)

type ProviderService struct {
	repo  *repo.Repos
	queue *queue.Queue
}

func NewProviderService(repo *repo.Repos, queue *queue.Queue) ProviderService {
	return ProviderService{repo: repo, queue: queue}
}

func (p *ProviderService) RequestCreateProvider(first_name, middle_name, second_name, email, phone string) error {
	providerInfo := dto.CreateProviderMessage{
		FirstName:  first_name,
		SecondName: second_name,
		Email:      email,
		Phone:      phone,
	}

	body, err := json.Marshal(providerInfo)
	if err != nil {
		return err
	}

	message := messages.Message{
		RoutingKey: "core.create_provider",
		Body:       body,
	}

	err = p.queue.SendMessage(message)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProviderService) CreateProvider(login, password, email string) error {
	pass := utils.GnerateHashPassword(password)
	err := p.repo.Authorization.CreateUser(login, pass, email, "PROVIDER")
	if err != nil {
		return err
	}

	return nil
}
