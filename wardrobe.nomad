job "wardrobe" {

  datacenters = ["dc1"]

  type = "service"

  group "wardrobe-stg-api" {

    count = 1

    network {

      port "api_port" {
        to = "8443"
      }
    }

    service {
      name = "wardrobe-stg-api"
      port = "8443"

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgres"
              local_bind_port  = 5432
            }
          }
        }
      }

      check {
        type = "http"
        port = "api_port"
        path = "/healthz"
        interval = "2s"
        timeout = "2s"
      }
    }

    task "stg-api" {
      driver = "docker"

      env {
        DEBUG = "true"
        PORT = ":8443"
        POSTGRES_CONNSTR = "postgres://wardrobe:wardrobe@127.0.0.1:5432/wardrobe?sslmode=disable"
        MINIO_USER = "minio"
        MINIO_PASSWORD = "secret"
        MINIO_ENDPOINT = "127.0.0.1:9000"
        FO_BUCKET = "fo-bucket"
      }

      config {
        image = "ghcr.io/gavrilaf/wardrobe-stg-api:0.0.1"
        ports = ["api_port"]
      }
    }
  }
}