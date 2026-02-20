# Documentation

This directory contains comprehensive documentation and diagrams for the Vapi AI Integration Backend.

## Files

### Main Documentation
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - Complete architectural documentation including:
  - System architecture overview
  - Provider abstraction strategy
  - Authentication flow
  - Call processing flow
  - Database schema
  - Deployment architecture
  - Testing strategy
  - Provider switching guide

### PlantUML Diagrams

All diagrams are in the `diagrams/` subdirectory and available in PlantUML format (.puml):

1. **[architecture.puml](./diagrams/architecture.puml)** - System architecture diagram showing hexagonal architecture with all layers (API, Application, Domain, Infrastructure)

2. **[provider-abstraction.puml](./diagrams/provider-abstraction.puml)** - Provider abstraction pattern showing how to easily switch from Vapi AI to other providers

3. **[auth-flow-sequence.puml](./diagrams/auth-flow-sequence.puml)** - Authentication and authorization flow including registration, login, token refresh

4. **[call-flow-sequence.puml](./diagrams/call-flow-sequence.puml)** - Complete call initiation and webhook processing flow

5. **[database-schema.puml](./diagrams/database-schema.puml)** - Database schema with all 6 tables, relationships, and indexes

6. **[deployment.puml](./diagrams/deployment.puml)** - Deployment architecture with Docker, load balancer, and external services

## Generating Diagrams

### Using PlantUML CLI

Install PlantUML:
```bash
# macOS
brew install plantuml

# Ubuntu/Debian
apt-get install plantuml

# Or download JAR from https://plantuml.com/download
```

Generate all diagrams:
```bash
cd diagrams
for file in *.puml; do
    plantuml "$file"
done
```

This will create PNG files for each diagram in the same directory.

Generate SVG instead:
```bash
for file in *.puml; do
    plantuml -tsvg "$file"
done
```

### Using Online Tools

You can also view and edit diagrams online:
- https://www.plantuml.com/plantuml/uml/
- https://plantuml-editor.kkeisuke.com/
- https://www.planttext.com/

Simply copy the contents of any `.puml` file and paste into the editor.

## Documentation Structure

```
docs/
├── README.md (this file)
├── ARCHITECTURE.md (comprehensive documentation)
└── diagrams/
    ├── architecture.puml
    ├── provider-abstraction.puml
    ├── auth-flow-sequence.puml
    ├── call-flow-sequence.puml
    ├── database-schema.puml
    └── deployment.puml
```

## Key Concepts

### Hexagonal Architecture
The system follows hexagonal architecture (ports & adapters) with:
- **Domain Layer**: Core business entities and interfaces
- **Application Layer**: Use cases and business logic services
- **Infrastructure Layer**: External integrations (database, providers)
- **API Layer**: HTTP handlers and middleware

### Provider Abstraction
The `VoiceProvider` interface allows easy switching between voice AI providers:
- Current implementation: Vapi AI
- Can add: Twilio, Vonage, custom providers
- **Zero changes** to business logic when switching

### Security
- JWT tokens (access + refresh)
- Bcrypt password hashing
- Webhook signature validation
- Input validation on all endpoints

### Testing
- Unit tests for services and entities
- Table-driven test approach
- Mock repositories and providers
- Current coverage: auth services and entities

## Quick Links

- [Main README](../README.md) - Project overview and setup
- [Quick Start Guide](../QUICKSTART.md) - Get started in 5 minutes
- [API Documentation](../API.md) - Complete API reference
- [Progress Tracking](../PROGRESS.md) - Implementation progress

## Contributing to Documentation

When updating documentation:
1. Keep PlantUML diagrams in sync with code changes
2. Update ARCHITECTURE.md for architectural changes
3. Regenerate PNG/SVG files after editing diagrams
4. Update this README if adding new documentation

## Tips

### PlantUML Syntax
- Use `@startuml` and `@enduml` to wrap diagrams
- Components, classes, sequences all supported
- Supports colors, notes, and formatting
- See https://plantuml.com/guide for full syntax

### Viewing in VS Code
Install PlantUML extension:
- Extension ID: `jebbs.plantuml`
- Preview: `Alt+D` or `Ctrl+Shift+P` → "PlantUML: Preview Current Diagram"

### Exporting for Presentations
Generate high-quality PNG for presentations:
```bash
plantuml -tpng -Sdpi=300 diagram.puml
```

## Questions?

For questions about the architecture or diagrams:
- Check [ARCHITECTURE.md](./ARCHITECTURE.md) first
- Review relevant diagram
- See [API.md](../API.md) for API specifics
- Check main [README.md](../README.md) for setup

---

Last Updated: 2024
