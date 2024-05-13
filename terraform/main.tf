provider "google" {
  credentials = file("${var.credentials_path}")
  project     = var.project_id
  region      = var.region
}

resource "google_cloud_run_service" "lux_service" {
  name     = "lux-service"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/${var.project_id}/realtime-analytics:latest"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_policy" "lux_service_iam" {
  location    = google_cloud_run_service.lux_service.location
  project     = google_cloud_run_service.lux_service.project
  service     = google_cloud_run_service.lux_service.name

  policy_data = <<EOF
{
  "bindings": [
    {
      "role": "roles/run.invoker",
      "members": [
        "allUsers"
      ]
    }
  ]
}
EOF
}
