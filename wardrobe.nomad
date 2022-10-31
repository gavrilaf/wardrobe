job "wardrobe" {

  datacenters = ["dc1"]

  type = "service"

  group "wardrobe" {
    network {
      port "wardrobe" {
        to = 8443
      }
    }

    count = 1

    task "storage" {
      driver = "docker"

      env {
        DEBUG: "true"
        PORT: ":8443"
        POSTGRES_CONNSTR: "postgres://wardrobe:wardrobe@127.0.0.1:5432/wardrobe?sslmode=disable"
        MINIO_USER: minio
        MINIO_PASSWORD: secret
        MINIO_ENDPOINT: "127.0.0.1:9000"
        FO_BUCKET: "fo-bucket"
      }

      config {
        image = "wardrobe:0.0.1"
        ports = ["wardrobe"]
      }
    }
  }
}