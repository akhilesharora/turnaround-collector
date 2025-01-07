# Turnaround Collector

## Overview

Turnaround Collector is a highly available, Go-based system designed for fetching images from multiple cameras and processing them through a target service.


## Prerequisites

- [Go 1.23](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

### Cloning the Repository

```bash
git clone https://github.com/akhilesharora/turnaround-collector.git
cd turnaround-collector
```

### Local Development

#### Build Services
```bash
make build
```

#### Run Unit Tests
```bash
make test
```

#### Run Integration Tests
```bash
make test-integration
```

### Docker Deployment

#### Build Docker Images
```bash
make docker-build
```

#### Start Services
```bash
make docker-up
```

#### View Logs
```bash
make docker-logs
```

#### Stop Services
```bash
make docker-down
```

## Configuration

The Collector service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `CAMERA_COUNT` | Number of camera replicas to poll | 3 |
| `MAX_CONCURRENT` | Maximum concurrent camera fetches | Half of camera count |
| `CAMERA_BASE_URL` | Base URL for camera service | `http://camera` |
| `TARGET_URL` | Full URL for image processing | `http://target:8080/image` |
| `POLL_INTERVAL` | Time between camera polls | 5 seconds |

## Architecture

### Components

1. **Camera Service**
    - Generates mock JPEG images
    - Endpoint: `/snap.jpg`
    - Simulates multiple camera sources

2. **Target Service**
    - Receives and processes images
    - Endpoint: `POST /image`
    - Logs image processing details

3. **Collector Service**
    - Polls cameras at configured intervals
    - Sends images to target service

### Workflow

1. Collector starts and configures camera polling
2. For each camera:
    - Fetch image from camera
    - Send image to target service
    - Log success or failure
3. Continue polling until context is canceled

## Error Handling

- Individual camera failures do not stop the entire system
- Errors are logged but do not interrupt other camera polling

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
