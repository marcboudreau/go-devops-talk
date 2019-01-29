terraform {
    required_version = "~> 0.11"
}

provider "google" {
    version = "1.19.0"
}

data "google_compute_image" "server_image" {
    family  = "go-devops-talk"
    project = "go-devops-talk"
}

resource "google_compute_instance" "docker" {
    name         = "docker"
    machine_type = "n1-standard-1"
    zone         = "us-central1-a"

    boot_disk {
        initialize_params {
            size  = "10"
            image = "${data.google_compute_image.server_image.self_link}"
        }
    }

    network_interface {
        network = "default"
        access_config {
        }
    }

    project = "go-devops-talk"
}
