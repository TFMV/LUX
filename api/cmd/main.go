package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/apache/iceberg-go/catalog"
	"github.com/apache/iceberg-go/table"
)

func main() {
	ctx := context.Background()
	projectID := "tfmv-371720"
	subID := "your-subscriber-id"

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	sub := client.Subscription(subID)

	// Initialize Iceberg catalog and table
	cat, tbl := initIcebergTable(ctx, "my_catalog_uri", "my_table_id")

	// Listen for messages
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Got message: %s\n", string(msg.Data))
		if err := processMessageIntoIceberg(ctx, tbl, msg.Data); err != nil {
			log.Printf("Failed to process message: %v", err)
			msg.Nack() // Nack message if processing fails
			return
		}
		msg.Ack()
	})

	if err != nil {
		log.Fatalf("Receive: %v", err)
	}
}

func initIcebergTable(ctx context.Context, uri string, tableID string) (catalog.Catalog, *table.Table) {
	// Initialize Iceberg catalog
	cat, err := catalog.NewRestCatalog("rest", uri)
	if err != nil {
		log.Fatalf("Failed to initialize catalog: %v", err)
	}

	// Load or create table
	tbl, err := cat.LoadTable(ctx, catalog.Identifier{tableID})
	if err != nil {
		log.Fatalf("Failed to load table: %v", err)
	}
	return cat, tbl
}

func processMessageIntoIceberg(ctx context.Context, tbl *table.Table, data []byte) error {
	// Unmarshal the JSON data into a generic map
	var record map[string]interface{}
	if err := json.Unmarshal(data, &record); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	// Retrieve the Iceberg schema
	schema := tbl.Schema()

	// Create a new generic Iceberg record
	icebergRecord, err := schema.NewRecord()
	if err != nil {
		return fmt.Errorf("failed to create a new record based on Iceberg schema: %w", err)
	}

	// Map data from JSON to Iceberg record based on the schema fields
	for _, field := range schema.Fields() {
		value, ok := record[field.Name()]
		if !ok {
			log.Printf("Warning: Missing data for field '%s'", field.Name())
			continue // or return an error if the field is required
		}

		if err := setIcebergRecordField(icebergRecord, field, value); err != nil {
			return fmt.Errorf("error setting field '%s': %w", field.Name(), err)
		}
	}

	// Append the record to the table
	if err := tbl.Append(ctx, icebergRecord); err != nil {
		return fmt.Errorf("failed to append record to Iceberg table: %w", err)
	}

	return nil
}

func setIcebergRecordField(rec *table.Record, field table.Field, value interface{}) error {
	switch field.Type() {
	case table.StringType:
		if str, ok := value.(string); ok {
			return rec.Set(field.Name(), str)
		}
	case table.IntegerType:
		if i, ok := value.(int); ok {
			return rec.Set(field.Name(), i)
		}
	case table.FloatType:
		if f, ok := value.(float64); ok {
			return rec.Set(field.Name(), f)
		}
	case table.BooleanType:
		if b, ok := value.(bool); ok {
			return rec.Set(field.Name(), b)
		}
	case table.TimestampType:
		if ts, ok := value.(string); ok {
			// Parse timestamp string into time.Time and convert to milliseconds
			parsedTime, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				return fmt.Errorf("invalid timestamp format for field '%s': %v", field.Name(), err)
			}
			return rec.Set(field.Name(), parsedTime.UnixNano()/1000000) // Convert to milliseconds
		}
	case table.DecimalType:
		if dec, ok := value.(string); ok {
			return rec.Set(field.Name(), dec)
		}
	case table.BinaryType:
		if bin, ok := value.([]byte); ok {
			return rec.Set(field.Name(), bin)
		}
	default:
		return fmt.Errorf("unsupported data type for field '%s'", field.Name())
	}

	return fmt.Errorf("expected %s value for field '%s', got %T", field.Type(), field.Name(), value)
}
