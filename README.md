# Overseerr Exporter    ![GitHub release (latest SemVer)](https://img.shields.io/github/v/tag/willfantom/overseerr-exporter?display_name=tag&label=%20&sort=semver)  ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/willfantom/overseerr-exporter/release.yml?label=%20&logo=github)

Export media request data from an [Oversseerr](https://overseerr.dev) instance to a [Prometheus](https://prometheus.io) instance.

---

## Usage

```bash
docker run --rm -p 9850:9850 ghcr.io/willfantom/overseerr-exporter:latest \
  "--overseerr.address=https://overseerr.example.com" \
  "--overseerr.api-key=examplesecretapikey"
```

---

### Dashboard

![example-dash](./grafana/dashboard.png)

---

## Exporter Metrics

Two main metric groups are exported: Requests and Users.

#### Requests

The requests on the Overseerr server are counted. Request counts have the following labels:

| Label            |                      Description                       | Configurable |
| :--------------- | :----------------------------------------------------: | -----------: |
| `request_status` |  The approval status for the requests (e.g. Approved)  |           no |
| `media_status`   | The media status for requested items (e.g. Available)  |           no |
| `media_type`     |       The category of request media (e.g. movie)       |           no |
| `is_4k`          |      Requested on a 4k tagged service (e.g. true)      |           no |
| `genre`          |       The main genre for a requested media item        |          yes |
| `company`        | The production company or network for a requested item |          yes |

> ⚠️  Collecting Genre/Company info can take a lot of time with large request quantities

#### Users

User request counts of an Overseerr server are collected with the following labels:

| Label   |          Description          | Configurable |
| :------ | :---------------------------: | -----------: |
| `email` | The email address of the user |           no |


## Configuration

| Flag                         |                 Description                 | Default    |
| :--------------------------- | :-----------------------------------------: | :--------- |
| `log`                        |   Sets the logging level for the exporter   | `fatal`    |
| `web.listen-address`         |  The address for the exporter to listen on  | `:9850`    |
| `web.telemetry-path`         |       The path to expose the metrics        | `/metrics` |
| `overseerr.address`          |      The URI of the Overseerr instance      |            |
| `overseerr.api-key`          | The admin API key of the Overseerr instance |            |
| `overseerr.locale`           |    The locale of the Overseerr instance     | `en`       |
| `overseerr.scrape.genres`    |   Collect genre information for requests    | `true`     |
| `overseerr.scrape.companies` |  Collect company information for requests   | `true`     |

You **must** provide the Overseerr address and API key!

---

## Build the Container

```bash
docker build --rm -f Dockerfile \
  --build-arg EXPORTER_VERSION=local \
  -t overseerr-exporter:latest .
```

---

### TODO

 - Improve dashboard (more graphs!)
 - Export version metrics
 - Include issue counters
