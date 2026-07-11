# Specs

This folder is the spec-first source of truth for the AI-native CI/CD platform.

## Order of Use

1. [OpenAPI contract](./openapi.yaml)
2. [Event architecture](../product/ai-cicd-platform/specs/event-architecture.md)
3. [Database schema](../product/ai-cicd-platform/specs/database-schema.md)
4. [Repository structure](../product/ai-cicd-platform/specs/repository-structure.md)
5. [Protobuf contract](../proto/ai_cicd_platform.proto)

## Development Rule

Implement new features from the contracts first. Update the spec before the service and UI code when behavior changes.
