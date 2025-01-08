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
# Public access is forbidden for now
# git clone https://github.com/akhilesharora/turnaround-collector.git
unzip turnaround-collector.zip
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

## Assumptions

- Camera Homogeneity: Assuming all cameras work the same way. Real-world, they might be different.
- Processing Speed: The target server is assumed to process images quickly. Eventually a queue system might be necessary to handle backpressure.
- Statelessness: The design is stateless now, but real use would need tracking processed images.
- Security: This implementation doesn't have any authentication or encryption.
- Image Size: The implementation assumes that images are of a reasonable size, huge images might need special handling.

## Error Handling

- Individual camera failures do not stop the entire system
- Errors are logged but do not interrupt other camera polling

## Improvements

- Retries and Circuit Breakers: This would make the system more reliable if things go wrong for a bit.
- Persistent Storage : Saving processed images or info could unlock new uses.
- Security: Authentication and encryption between services is key for real use.
- Tracing: Distributed tracing would make debugging and performance analysis easier.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

