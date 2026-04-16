# Faro: A Middleware for Scalable Medical Imaging Data Warehouse Construction and Deduplication

## Overview

Faro is a high-performance near-duplicate detection engine designed for the construction of medical imaging data warehouses. Built in Go, it translates the research framework established in the MediCurator project into a scalable, cloud-native middleware.

Faro integrates data from heterogeneous sources, detects duplicates across multiple dimensions, and consolidates records into a unified warehouse.

## Research Basis

The core methodologies include:
- Similarity matrices for textual medical records.
- Hierarchical metadata comparison (Study/Series/Image) for binary medical archives.
- Distributed execution patterns for high-throughput processing.

## Features

- **Clinical Normalization Layer**: Automatically handles medical acronyms (e.g., ASA, HCTZ) and standardizes measurement units for high-precision matching.
- **Hierarchical Metadata Comparison**: Identifies duplicates in massive medical image archives by analyzing Study and Series hierarchies without processing raw binary files.
- **Hybrid Storage Layer**: Combines thread-safe in-memory processing with persistent Key-Value storage (BadgerDB) for resilient and scalable operations.
- **Go-Native Concurrency**: Leverages goroutines and channels for efficient, parallel duplicate scanning.

## Architecture

Faro is structured as modular Go packages:
- `internal/pkg/detect`: Detection algorithms and normalization logic.
- `internal/pkg/storage`: Storage implementations for in-memory and persistent KV stores.
- `internal/pkg/types`: Core medical data and imaging hierarchy models.

## Installation

### Prerequisites (Linux)

If the Go and Docker toolchains are missing in your environment, you can install them via the following commands:

```bash
# Install Go 1.21+
sudo apt update
sudo apt install golang-go

# Install Docker and Docker Compose
sudo apt install docker.io docker-compose
```

### Setup
```bash
# Clone the repository
git clone [repository-url]
cd Faro

# Initialize dependencies
go mod tidy
```

## Execution

### Option 1: Single-Node (1-Click Demo)
This mode runs the Faro Master engine with an in-memory database and an embedded SQLite demonstration. It is ideal for verifying the core discovery logic and normalization.

```bash
go run cmd/faro/main.go
```
The **Resolution Dashboard** will be available at: `http://localhost:8080`

### Option 2: Distributed Mode (Resilient Cluster)
This mode uses NATS JetStream for persistent job distribution across multiple worker nodes.

### Running in Distributed Mode (Multi-Node)
To test the Faro cluster locally using Docker Compose:

1. **Build and Scale the Cluster**:
   ```bash
   docker-compose up -d --build
   docker-compose scale worker=3
   ```

2. **Run the Master Node**:
   The `main.go` entry point will detect the NATS server and begin distributing jobs to the pool of workers.
   ```bash
   go run cmd/faro/main.go
   ```

3. **Monitor Logs**:
   ```bash
   docker-compose logs -f worker
   ```

## Architecture
Faro is designed as a distributed middleware:
- **Master Node**: Fetches data from sources and orchestrates comparison jobs.
- **Worker Node**: Stateless nodes that perform computationally expensive similarity calculations.
- **NATS JetStream**: Provides a persistent, resilient messaging backbone for job distribution.
