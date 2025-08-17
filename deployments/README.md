# Deployments

This directory contains deployment-related configuration files.

## Files

- **Dockerfile** - Multi-stage Docker image definition
- **docker-compose.yml** - Docker Compose configuration

## Quick Commands

```bash
# Build and run with Makefile
make docker-build
make docker-up

# Manual Docker commands
docker build -t loan-service -f deployments/Dockerfile .
docker run -p 8080:8080 loan-service
docker-compose -f deployments/docker-compose.yml up -d
```

For complete deployment guide, see [docs/README.md](../docs/README.md#deployment-guide).
