package ingester

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"moviemicroservice.com/src/modules/ratings/pkg/models"
)

type Ingester struct {
	consumer *kafka.Consumer
	topic    string
}

func New(addr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": addr,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &Ingester{consumer, topic}, err
}

func (i *Ingester) Ingest(ctx context.Context) (chan models.RatingEvent, error) {
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan models.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
				return
			default:
			}

			//indefinite wait
			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Println("consumer read message error: " + err.Error())
				continue
			}

			var event models.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Println("json decode error: " + err.Error())
				continue
			}

			ch <- event
		}
	}()

	return ch, nil
}
