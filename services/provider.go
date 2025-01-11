package services

import (
	"auth/queue"
	"auth/queue/messages"
	"auth/queue/messages/dto"
	"auth/repo"
	"encoding/json"
)

type ProviderService struct {
	repo  *repo.Repos
	queue *queue.Queue
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
	return nil
}
