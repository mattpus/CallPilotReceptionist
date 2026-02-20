# CallPilotReceptionist

A scalable Go-based backend for integrating Vapi AI voice assistants with small businesses. This system provides call handling, appointment scheduling, and a comprehensive dashboard for monitoring AI voice interactions.

## 🎯 Project Status

**Status**: ✅ **Live and Working!** 

- ✅ Complete backend implementation (37 Go files, 4,817 lines)
- ✅ 20+ REST API endpoints
- ✅ Provider abstraction for easy switching
- ✅ JWT authentication
- ✅ Supabase database integration
- ✅ Webhook processing (Vapi events)
- ✅ Cloudflare Tunnel support
- ✅ Docker deployment ready
- ✅ Comprehensive documentation (2,709 lines)
- ✅ 6 PlantUML architecture diagrams
- ✅ Unit tests for core components

**Test Coverage**: ~65% (target: 80%+)

**Live Features:**
- 📞 Real phone number integration via Vapi
- 🔄 Webhook event processing
- 💾 Call logging to Supabase database
- 🔐 Secure credential management

## 📚 Documentation

### Quick Links
- **[Quick Start Guide](QUICKSTART.md)** - Get started in 5 minutes
- **[API Documentation](API.md)** - Complete API reference with examples
- **[Architecture Documentation](docs/ARCHITECTURE.md)** - Deep dive into system design
- **[Testing Guide](docs/TESTING.md)** - Testing strategy and coverage
- **[Progress Tracking](PROGRESS.md)** - Implementation status

### Architecture Diagrams
All diagrams available in PlantUML format in `docs/diagrams/`:
- **System Architecture** - Hexagonal architecture overview
- **Provider Abstraction** - How to switch voice AI providers
- **Authentication Flow** - Registration, login, token refresh
- **Call Processing Flow** - End-to-end call handling
- **Database Schema** - 6 tables with relationships
- **Deployment Architecture** - Docker multi-container setup

## 🏗️ Architecture

The project follows **Hexagonal Architecture** (Ports and Adapters) principles:

- **Domain Layer**: Business entities and core logic
- **Application Layer**: Use cases and services
- **Infrastructure Layer**: External integrations (Vapi AI, Supabase)
- **API Layer**: HTTP handlers and middleware

### Key Features

- ✅ **Provider Abstraction**: Easy switching between voice AI providers (zero business logic changes)
- ✅ **Comprehensive Testing**: Unit tests, mocks, and integration tests
- ✅ **JWT Authentication**: Secure business user authentication (access + refresh tokens)
- ✅ **Call Monitoring**: Track all calls and interactions in real-time
- ✅ **Appointment Extraction**: Automatically parse appointment requests from calls
- ✅ **Analytics Dashboard**: Business insights, call volume, and trends
- ✅ **Webhook Processing**: Secure webhook handling with signature validation
- ✅ **Async Transcript Fetch**: Non-blocking transcript retrieval after calls

## 📁 Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/
│   │   └── handlers/    # HTTP request handlers
│   ├── application/
│   │   ├── dto/         # Data transfer objects
│   │   └── services/    # Business logic services
│   ├── domain/
│   │   ├── entities/    # Business entities
│   │   ├── errors/      # Domain errors
│   │   └── providers/   # Provider interfaces
│   └── infrastructure/
│       ├── database/    # Database repositories
│       ├── http/
│       │   └── middleware/ # HTTP middleware
│       └── providers/
│           └── vapi/    # Vapi AI implementation
├── pkg/
│   ├── config/          # Configuration management
│   ├── logger/          # Logging utilities
│   └── utils/           # Shared utilities
├── migrations/          # Database migrations
└── docker-compose.yml   # Docker orchestration
```

## 🚀 Quick Start

### Prerequisites

- Docker Desktop (for running the application)
- Supabase account (free tier works)
- Vapi AI account and API key
- Cloudflare Tunnel (cloudflared) for webhook exposure

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/mattpus/CallPilotReceptionist.git
   cd CallPilotReceptionist
   ```

2. **Set up Supabase Database**
   - Sign up at https://supabase.com
   - Create a new project
   - Go to Settings → Database
   - Copy the **Connection Pooling** connection string (URI format)
   - Run the migration SQL from `migrations/001_initial_schema.up.sql` in SQL Editor

3. **Create environment file**
   ```bash
   cp .env.example .env
   ```

4. **Configure environment variables** (edit `.env`):
   ```bash
   # Database (from Supabase)
   DATABASE_URL=postgresql://postgres.xxxxx:password@aws-x-region.pooler.supabase.com:5432/postgres
   
   # JWT Secret (generate with: openssl rand -base64 32)
   JWT_SECRET_KEY=your-generated-secret-here
   
   # Vapi AI
   VAPI_API_KEY=your-vapi-api-key
   VAPI_WEBHOOK_URL=https://your-tunnel-url.trycloudflare.com/api/v1/webhooks/vapi
   ```

5. **Start the application with Docker**
   ```bash
   docker-compose up --build -d
   ```

6. **Set up Cloudflare Tunnel** (for webhook access)
   ```bash
   # Install cloudflared
   brew install cloudflared
   
   # Run tunnel (keep this running)
   cloudflared tunnel --url http://localhost:8080
   ```
   
   Copy the HTTPS URL and update `VAPI_WEBHOOK_URL` in `.env`, then restart:
   ```bash
   docker-compose restart app
   ```

7. **Verify the server is running**
   ```bash
   curl http://localhost:8080/health
   # Should return: {"status":"ok"}
   ```

8. **Configure Vapi Assistant**
   - Log into https://dashboard.vapi.ai
   - Create a new Assistant
   - Add your webhook URL in Server Messages/Server URL
   - Buy a phone number and link to your assistant
   - Call the number to test!

## 🔧 Configuration

All configuration is managed through environment variables (see `.env.example` for template).

### Required Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | Supabase connection string (pooler) | `postgresql://postgres.xxx:pass@aws-x.pooler.supabase.com:5432/postgres` |
| `JWT_SECRET_KEY` | Secret for JWT signing | Generate with `openssl rand -base64 32` |
| `VAPI_API_KEY` | Vapi AI API key | From dashboard.vapi.ai |
| `VAPI_WEBHOOK_URL` | Public webhook endpoint | `https://xxx.trycloudflare.com/api/v1/webhooks/vapi` |

### Optional Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` |
| `DB_MAX_OPEN_CONNS` | Max database connections | `25` |

### Webhook URL Setup

You need a publicly accessible URL for Vapi to send webhooks. Options:

**Option 1: Cloudflare Tunnel (Recommended for Development)**
```bash
brew install cloudflared
cloudflared tunnel --url http://localhost:8080
# Copy the https://xxx.trycloudflare.com URL
```

**Option 2: Deploy to Cloud (Production)**
- Railway: `railway up`
- Render: Connect GitHub repo
- Fly.io: `fly deploy`

## 🧪 Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run integration tests:
```bash
go test -tags=integration ./...
```

### Testing Webhooks

Test that your webhook endpoint is accessible:
```bash
curl -X POST https://your-tunnel-url.trycloudflare.com/api/v1/webhooks/vapi \
  -H "Content-Type: application/json" \
  -d '{"test": "ping"}'
```

You should see a 401 response (expected - validates webhook is working but rejects unsigned requests).

## 📚 API Documentation

Once the server is running, API documentation will be available at:
- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI Spec: `http://localhost:8080/openapi.json`

### Key Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register new business
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh JWT token

#### Calls
- `GET /api/v1/calls` - List all calls
- `GET /api/v1/calls/:id` - Get call details
- `POST /api/v1/calls` - Initiate a call
- `POST /api/v1/webhooks/vapi` - Vapi AI webhook

#### Analytics
- `GET /api/v1/analytics/overview` - Dashboard stats
- `GET /api/v1/analytics/calls` - Call volume metrics

## 🔄 Switching Voice AI Providers

The system is designed for easy provider switching:

1. **Implement the `VoiceProvider` interface** in `internal/domain/providers/`
2. **Create provider implementation** in `internal/infrastructure/providers/[provider-name]/`
3. **Update provider factory** in configuration
4. **No changes needed** to services or handlers

Example:
```go
type VoiceProvider interface {
    InitiateCall(ctx context.Context, req CallRequest) (*CallSession, error)
    HandleWebhook(ctx context.Context, payload []byte) (*CallEvent, error)
    GetCallDetails(ctx context.Context, callID string) (*CallDetails, error)
    GetTranscript(ctx context.Context, callID string) (*Transcript, error)
}
```

## 🗄️ Database Schema

The application uses PostgreSQL via Supabase with the following main tables:
- `businesses` - Business information and settings
- `users` - User accounts with JWT authentication
- `calls` - Call records with status and metadata
- `interactions` - Call interactions and events
- `transcripts` - Full call transcripts
- `appointments` - Extracted appointment information

### Running Migrations

Migrations are in the `migrations/` directory.

**For Supabase:**
1. Go to SQL Editor in Supabase dashboard
2. Copy content from `migrations/001_initial_schema.up.sql`
3. Run the query

**For local PostgreSQL:**
```bash
psql $DATABASE_URL -f migrations/001_initial_schema.up.sql
```

## 🚢 Deployment

### Docker (Recommended)

Build and run with Docker:
```bash
docker-compose up --build -d
```

Check logs:
```bash
docker-compose logs -f app
```

Stop:
```bash
docker-compose down
```

### Cloud Platforms

**Railway**
```bash
railway init
railway up
```

**Render**
1. Connect your GitHub repository
2. Add environment variables in dashboard
3. Render will auto-deploy on push

**Fly.io**
```bash
fly launch
fly deploy
```

### Production Checklist

- [ ] Set `ENVIRONMENT=production` in environment
- [ ] Use strong JWT secret (32+ characters)
- [ ] Configure proper `DATABASE_URL` with connection pooling
- [ ] Set up monitoring/logging
- [ ] Configure CORS if needed
- [ ] Set up SSL/TLS (automatic with Cloudflare/Railway/Render)
- [ ] Configure webhook URL with permanent domain

## 📝 Development

### Adding a New Feature

1. Define domain entities in `internal/domain/entities/`
2. Create service in `internal/application/services/`
3. Implement repository in `internal/infrastructure/database/`
4. Add HTTP handlers in `internal/api/handlers/`
5. Write tests for all layers

### Code Style

Follow standard Go conventions:
- Use `gofmt` for formatting
- Run `golint` for linting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

[Your License Here]

## 🆘 Support

For issues and questions:
- GitHub Issues: https://github.com/mattpus/CallPilotReceptionist/issues
- Documentation: See `/docs` folder
- Security: See `SECURITY.md`

## 🐛 Troubleshooting

### Webhook not receiving events
- Check Cloudflare tunnel is running: Look for "Your quick Tunnel has been created"
- Verify webhook URL in Vapi dashboard matches your tunnel URL
- Test endpoint: `curl https://your-url.trycloudflare.com/health`
- Check logs: `docker-compose logs -f app`

### Database connection errors
- Verify `DATABASE_URL` is the **Connection Pooling** URL from Supabase
- Format: `postgresql://postgres.xxxxx:password@aws-x.pooler.supabase.com:5432/postgres`
- Check migrations have been run in Supabase SQL Editor
- Test connection in Supabase dashboard

### Docker build fails
- Ensure Docker Desktop is running
- Clear cache: `docker-compose down && docker system prune -a`
- Rebuild: `docker-compose up --build`

### "Call not found" in webhook logs
- This is normal if call was initiated directly from Vapi dashboard
- Calls initiated via API will be tracked properly
- Webhooks are still being received and processed correctly

### Environment variables not loading
- Restart Docker after changing `.env`: `docker-compose restart app`
- Verify `.env` exists and is not `.env.example`
- Check for syntax errors in `.env` file

## 🗺️ Roadmap

### Completed ✅
- Backend server with REST API
- Supabase database integration
- Vapi webhook processing
- JWT authentication
- Docker deployment
- Cloudflare tunnel support

### In Progress 🚧
- Increase test coverage to 80%+
- Frontend dashboard
- Additional voice provider support

### Planned 📋
- Real-time call monitoring dashboard
- Advanced analytics
- Multi-language support
- Custom AI prompts per business
- Appointment scheduling integration (Google Calendar, etc.)

## 📝 License

[Your License Here]

## 🙏 Acknowledgments

- Built with Vapi AI for voice capabilities
- Supabase for database infrastructure
- Cloudflare for tunnel services
