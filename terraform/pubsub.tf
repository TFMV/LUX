resource "google_pubsub_topic" "iot_topic" {
  name = "iot-data-topic"
}

resource "google_pubsub_subscription" "iot_subscription" {
  name  = "iot-subscription"
  topic = google_pubsub_topic.iot_topic.name

  // Acknowledgement deadline (the time after which the message will be redelivered if not acknowledged)
  ack_deadline_seconds = 20

  // Message retention duration (the time that Pub/Sub retains the message without acknowledgment)
  message_retention_duration = "86400s"  // 24 hours

  // Expiration policy (how long a subscription should exist when no messages are being sent)
  expiration_policy {
    ttl = "2678400s" // 31 days
  }

  // Retry policy (how Pub/Sub retries delivering a message in case of failure)
  retry_policy {
    minimum_backoff = "10s"
    maximum_backoff = "600s"
  }

  // Dead letter policy (configuring a dead letter topic for undeliverable messages)
  dead_letter_policy {
    dead_letter_topic = google_pubsub_topic.dead_letter_topic.id
    max_delivery_attempts = 5
  }
}

resource "google_pubsub_topic" "dead_letter_topic" {
  name = "iot-dead-letter-topic"
}
