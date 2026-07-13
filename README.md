# Kinetic

A distributed DAG execution runtime for fault-tolerant workflow orchestration, resource-aware scheduling, and AI-native pipelines.

<img width="1595" height="595" alt="Screenshot From 2026-07-13 19-11-50" src="https://github.com/user-attachments/assets/bb5779ce-74ca-4cba-865f-eb9cff1ef743" />

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
