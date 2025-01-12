package consumer

import (
	"auth/queue/messages"
	"auth/queue/messages/dto"
	"auth/services"
	"encoding/json"
	"log"
)

type Consumer struct {
	services *services.Services
}

func NewConsumer(services *services.Services) *Consumer {
	return &Consumer{services: services}
}

func (c *Consumer) Consume(qMessage []byte) {
	var message messages.Message

	err := json.Unmarshal(qMessage, &message)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if message.RoutingKey == "auth.created_provider" {
		var msg = dto.CreatedProviderMessage{}
		err = json.Unmarshal(message.Body, &msg)
		if err != nil {
			log.Println(err.Error())
			return
		}
		c.services.ProviderService.CreateProvider(msg.Login, msg.Password, msg.Email)
	}
}
