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
  - job_name: 'Encoder_Scrape'
    static_configs:
      - targets: ['direkt_exporter:9110']
        labels:
          name: 'Encoder'
    scrape_interval: 600s
    scrape_timeout: 15s
    params:
      serial:
        - D01234
    metrics_path: /probe
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: direkt_exporter:9110
```

### Building

```
docker build -t direkt:latest .
```