package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
)

type Record struct {
	Data string `json:"data"`
}

func main() {
	ctx := context.Background()
	projectID := "tfmv-371720"
	subID := "your-subscriber-id"

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	sub := client.Subscription(subID)
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		// Convert message to JSON
		record := Record{Data: string(msg.Data)}
		recordBytes, err := json.Marshal(record)
		if err != nil {
			log.Printf("Error marshaling record: %v", err)
			msg.Nack()
			return
		}

		// Send data to the middleware Iceberg service
		resp, err := http.Post("http://localhost:8080/append", "application/json", bytes.NewReader(recordBytes))
		if err != nil {
			log.Printf("Failed to send data to Iceberg service: %v", err)
			msg.Nack()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to append data to Iceberg table: %s", resp.Status)
			msg.Nack()
		} else {
			msg.Ack()
		}
	})

	if err != nil {
		log.Fatalf("Receive: %v", err)
	}
}
