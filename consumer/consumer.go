package consumer

import (
	"auth/queue/messages"
	"auth/queue/messages/dto"
	"auth/services"
	"encoding/json"
	"log"
)

type Consumer struct {
	providerService services.ProviderService
}

func (c *Consumer) Consume(qMessage []byte) {
	var message messages.Message

	err := json.Unmarshal(qMessage, &message)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if message.RoutingKey == "auth.create_provider" {
		var msg = dto.CreatedProviderMessage{}
		err = json.Unmarshal(message.Body, &msg)
		if err != nil {
			log.Println(err.Error())
			return
		}
		c.providerService.CreateProvider(msg.Login, msg.Password, msg.Email)
	}
}
