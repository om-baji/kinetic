# Kinetic

A distributed DAG execution runtime for fault-tolerant workflow orchestration, resource-aware scheduling, and AI-native pipelines.

## Features

- Distributed DAG execution
- Dependency-aware scheduling
- Fault-tolerant worker runtime
- Resource-aware task placement
- Checkpointing and automatic recovery
- Retry policies with configurable backoff
- Event-driven workflow orchestration
- Real-time observability with logs and metrics

## Tech Stack

- **Backend:** Go
- **API:** gRPC, REST
- **Database:** PostgreSQL
- **Queue/Event Bus:** In-memory (Redis/Kafka planned)
- **Observability:** Prometheus, OpenTelemetry, Grafana

## Architecture

The overall system architecture is shown below.
