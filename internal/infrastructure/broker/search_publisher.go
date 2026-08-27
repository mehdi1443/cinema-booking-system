package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/broker"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	SearchExchange   = "search_events_exchange"
	SyncRoutingKey   = "search.sync"
	DeleteRoutingKey = "search.delete"
)

type searchPublisher struct {
	channel *amqp.Channel
}

func NewSearchEventPublisher(ch *amqp.Channel) (broker.SearchEventPublisher, error) {
	err := ch.ExchangeDeclare(
		SearchExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &searchPublisher{channel: ch}, nil
}

func (p *searchPublisher) PublishSyncEvent(ctx context.Context, doc entity.SearchDocument) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		SearchExchange,
		SyncRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish sync event: %w", err)
	}
	return nil
}

func (p *searchPublisher) PublishDeleteEvent(ctx context.Context, documentID string) error {
	payload := map[string]string{"id": documentID}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal delete payload: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		SearchExchange,
		DeleteRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish delete event: %w", err)
	}
	return nil
}
