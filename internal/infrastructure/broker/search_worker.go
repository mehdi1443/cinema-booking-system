package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"github.com/mehdi1443/cinema-booking-system/internal/domain/repository/elastic"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	SyncQueue   = "search_sync_queue"
	DeleteQueue = "search_delete_queue"
)

type SearchWorker struct {
	channel    *amqp.Channel
	searchRepo elastic.SearchRepository
	log        domainlogger.Logger
}

func NewSearchWorker(
	ch *amqp.Channel,
	searchRepo elastic.SearchRepository,
	log domainlogger.Logger,
) (*SearchWorker, error) {

	err := ch.ExchangeDeclare(SearchExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange in worker: %w", err)
	}

	_, err = ch.QueueDeclare(SyncQueue, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare sync queue: %w", err)
	}
	err = ch.QueueBind(SyncQueue, SyncRoutingKey, SearchExchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind sync queue: %w", err)
	}

	_, err = ch.QueueDeclare(DeleteQueue, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare delete queue: %w", err)
	}
	err = ch.QueueBind(DeleteQueue, DeleteRoutingKey, SearchExchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind delete queue: %w", err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &SearchWorker{
		channel:    ch,
		searchRepo: searchRepo,
		log:        log,
	}, nil
}

func (w *SearchWorker) Start(ctx context.Context) error {
	syncMsgs, err := w.channel.Consume(SyncQueue, "sync_consumer", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register sync consumer: %w", err)
	}

	deleteMsgs, err := w.channel.Consume(DeleteQueue, "delete_consumer", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register delete consumer: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	w.log.Info(ctx, "Search Worker started successfully, waiting for messages...")

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				w.log.Info(ctx, "Shutting down sync consumer...")
				return
			case msg, ok := <-syncMsgs:
				if !ok {
					return
				}
				w.processSyncMsg(ctx, msg)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				w.log.Info(ctx, "Shutting down delete consumer...")
				return
			case msg, ok := <-deleteMsgs:
				if !ok {
					return
				}
				w.processDeleteMsg(ctx, msg)
			}
		}
	}()

	wg.Wait()
	return nil
}

func (w *SearchWorker) processSyncMsg(ctx context.Context, msg amqp.Delivery) {
	var doc entity.SearchDocument
	if err := json.Unmarshal(msg.Body, &doc); err != nil {
		w.log.Error(ctx, "Failed to unmarshal sync message", domainlogger.Error(err))
		msg.Nack(false, false)
		return
	}

	err := w.searchRepo.Index(ctx, doc)
	if err != nil {
		w.log.Error(ctx, "Elasticsearch Index failed", domainlogger.Error(err), domainlogger.String("doc_id", doc.ID))
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	w.log.Info(ctx, "Document indexed in Elasticsearch successfully", domainlogger.String("doc_id", doc.ID))
}

func (w *SearchWorker) processDeleteMsg(ctx context.Context, msg amqp.Delivery) {
	var payload map[string]string
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		w.log.Error(ctx, "Failed to unmarshal delete message", domainlogger.Error(err))
		msg.Nack(false, false)
		return
	}

	docID, ok := payload["id"]
	if !ok {
		w.log.Error(ctx, "Delete message missing 'id' field")
		msg.Nack(false, false)
		return
	}

	err := w.searchRepo.Delete(ctx, docID)
	if err != nil {
		w.log.Error(ctx, "Elasticsearch Delete failed", domainlogger.Error(err), domainlogger.String("doc_id", docID))
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	w.log.Info(ctx, "Document deleted from Elasticsearch successfully", domainlogger.String("doc_id", docID))
}
