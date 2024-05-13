resource "google_pubsub_subscription" "iot_subscription" {
  name  = "iot-subscription"
  topic = google_pubsub_topic.iot_topic.name

  push_config {
    push_endpoint = google_cloud_run_service.lux_service.status[0].url

    oidc_token {
      service_account_email = google_service_account.lux_invoker.email
    }
  }
}

resource "google_service_account" "lux_invoker" {
  account_id   = "lux-invoker"
  display_name = "LUX Invoker Account"
}

resource "google_project_iam_member" "invoker_member" {
  project = google_project.default.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.lux_invoker.email}"
}
