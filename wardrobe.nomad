job "wardrobe" {

  datacenters = ["dc1"]

  type = "service"

  group "storage" {
    network {
      port "storage" {
        to = 4563
      }
    }

    count = 1

    task "storage" {
      driver = "docker"

      config {
        image = "wardrobe:0.0.1"
        ports = ["storage"]
      }
    }
  }
}