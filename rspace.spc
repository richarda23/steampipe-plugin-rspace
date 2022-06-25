connection "rspace" {
    plugin    = "local/rspace"
    options "connection" {
      cache     = false # true, false
      cache_ttl = 300  # expiration (TTL) in seconds
    }
}