# Yandex Cloud Object Storage Service

A Go-based microservice that provides a REST API for interacting with Yandex Cloud Object Storage. This service simplifies common operations like listing objects and generating pre-signed URLs for secure object access.

## Features

- List objects in a bucket with prefix filtering
- Generate pre-signed URLs for temporary object access
- Environment-based configuration
- Health check endpoint
- Built with Gin web framework

## Prerequisites

- Go 1.21 or later
- Yandex Cloud account with Object Storage enabled
- S3 compatible storage access credentials

## Configuration

Copy the example environment file and configure your settings:

```bash
cp .env.example .env
```

Environment variables:

- `SERVER_PORT`: HTTP server port (default: 8081)
- `S3_ENDPOINT`: Yandex Object Storage endpoint (default: storage.yandexcloud.net)
- `S3_REGION`: Storage region (default: ru-central1)
- `S3_BUCKET_NAME`: Your bucket name
- `S3_ACCESS_KEY_ID`: Your access key ID
- `S3_SECRET_ACCESS_KEY`: Your secret access key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/marketconnect/yc-object-storage-service.git
cd yc-object-storage-service
```

2. Install dependencies:
```bash
go mod download
```

3. Build the service:
```bash
go build
```

## Usage

Start the service:
```bash
./storage-service
```

### API Endpoints

#### List Objects
```http
GET /api/v1/list?prefix=path/to/folder
```
Query Parameters:
- `prefix`: (Required) The prefix to filter objects by

Response:
```json
{
    "objects": [
        "path/to/folder/file1.txt",
        "path/to/folder/file2.jpg"
    ]
}
```

#### Generate Pre-signed URL
```http
GET /api/v1/generate-url?objectKey=path/to/file&expires=3600
```
Query Parameters:
- `objectKey`: (Required) The key of the object to generate URL for
- `expires`: (Optional) URL expiration time in seconds (default: 3600)

Response:
```json
{
    "url": "https://storage.yandexcloud.net/..."
}
```

#### Health Check
```http
GET /health
```
Response:
```json
{
    "status": "ok"
}
```

## Error Handling

The service returns appropriate HTTP status codes and error messages:

- 400 Bad Request: Missing or invalid parameters
- 500 Internal Server Error: Server-side errors with details in the response

Example error response:
```json
{
    "error": "failed to list objects",
    "details": "error details here"
}
```

## Development

The project structure follows standard Go project layout:

```
.
├── api/            # HTTP handlers
├── config/         # Configuration management
├── s3/            # S3 client implementation
├── main.go        # Application entry point
├── go.mod         # Go modules file
└── .env.example   # Example environment configuration
```

## License

[Add your license information here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
