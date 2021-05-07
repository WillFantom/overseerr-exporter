# Overseerr Exporter

Export media request data from an [Oversseerr](https://overseerr.dev) instance to a [Prometheus](https://prometheus.io) instance.

### Exporter Metrics

  - Requests:
    - By media type (tv, movie, etc...)
    - By status (available, approved, pending, declined)
    - By 4k (regular or 4k)

  - Users:
    - With request count

  - Genres
    - Request genre counts
    > This scrape can take a large amount of time if you have a lot of requests!

### Exporter Config

  - Default port: `9850`
  
##### Required Values:
  * `overseerr.address`: the URI of your Overseerr instance
  * `overseerr.api-key`: the admin API key of your Overseerr instance

## Usage

Using the binary:

- Download the appropiate binary version from the GitHub releases page

- ```bash
  overseerr-exporter \
    --overseerr.address=https://overseerr.example.com \
    --overseerr.api-key=examplesecretapikey
  ```

Using Docker:

```bash
docker run --rm -p 9850:9850 ghcr.io/willfantom/overseerr-exporter:latest \
  "--overseerr.address=https://overseerr.example.com" \
  "--overseerr.api-key=examplesecretapikey"
```

## Build the Container

```bash
docker build --rm -f Dockerfile --build-arg EXPORTER_VERSION=local \
  -t overseerr-exporter:latest .
```
