job "wardrobe-external" {

  datacenters = ["dc1"]

  group "wardrobe-external" {
    count = 1
    network {
      mode = "bridge"
    }

    service {
      name = "wardrobe-external"

      connect {
        gateway {
          proxy {}
          terminating {
            service {
              name = "postgres"
            }
          }
        }
      }
    }
  }
}