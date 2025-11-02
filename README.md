# Prometheus Direkt Exporter

> Warning: Stability is not guaranteed until v1 release

## Usage

### Docker
```
docker run --rm \
  -p 9110/tcp \
  --name direkt_exporter \
  IMAGE_TBC:latest
  ```
- If using authentication, set the DIREKT_USERNAME and DIREKT_PASSWORD environment variables when running the docker image.

Healthcheck is available on /-/healthy

### Prometheus Config
 
 Example config
 ```
 scrape_configs:
  - job_name: 'blackbox'
    metrics_path: /probe
    static_configs:
      - targets:
        - D01111   # Target serials
        - D01234
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9110  # The exporter's real hostname:port.
  - job_name: 'direkt_exporter'  # collect blackbox exporter's operational metrics.
    static_configs:
      - targets: ['127.0.0.1:9110']
```