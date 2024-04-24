package producer

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"moviemicroservice.com/services/ratings/pkg/models"
)

func main() {
	fmt.Println("creating kafka producer")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	const filename = "ratings_data.json"
	fmt.Println("reading sample events from file" + filename)

	events, err := readEvents(filename)
	if err != nil {
		panic(err)
	}

	const topic = "ratings"
	if err := produceEvents(topic, producer, events); err != nil {
		panic(err)
	}

	const timeout = 10 * time.Second

	fmt.Println("Waiting for " + timeout.String() + " until all events are produced")

	producer.Flush(int(timeout.Milliseconds()))
}

func readEvents(filename string) ([]models.RatingEvent, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var events []models.RatingEvent
	if err := json.NewDecoder(f).Decode(&events); err != nil {
		return nil, err
	}

	return events, nil
}

func produceEvents(topic string, producer *kafka.Producer, events []models.RatingEvent) error {
	for _, event := range events {
		encoded, err := json.Marshal(event)
		if err != nil {
			return err
		}

		if err := producer.Produce(&kafka.Message{Value: []byte(encoded), TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}}, nil); err != nil {
			return err
		}
	}

	return nil
}
