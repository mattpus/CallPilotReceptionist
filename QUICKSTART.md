# Quick Start Guide

Get your Vapi AI Integration backend running in 5 minutes!

## Prerequisites

- Go 1.21 or later
- Supabase account (or PostgreSQL database)
- Vapi AI account and API key

## Step 1: Clone and Configure

```bash
# Navigate to project directory
cd vapiAIIntegration

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env
```

Required environment variables:
- `DATABASE_URL` - Your Supabase/PostgreSQL connection string
- `JWT_SECRET_KEY` - Generate with: `openssl rand -base64 32`
- `VAPI_API_KEY` - Your Vapi AI API key from dashboard
- `VAPI_WEBHOOK_URL` - Your public webhook URL (for production)

## Step 2: Set Up Database

Connect to your Supabase database and run the migration:

```bash
# Using psql
psql $DATABASE_URL < migrations/001_initial_schema.up.sql

# Or using Supabase dashboard:
# - Go to SQL Editor
# - Paste contents of migrations/001_initial_schema.up.sql
# - Run query
```

This creates 6 tables with indexes:
- businesses
- users
- calls
- interactions
- transcripts
- appointments

## Step 3: Install Dependencies

```bash
go mod download
```

## Step 4: Build and Run

```bash
# Build
go build -o server cmd/server/main.go

# Run
./server

# Or run directly
go run cmd/server/main.go
```

You should see:
```
{"level":"info","time":"...","message":"Starting CallPilotReceptionist Server","environment":"development","port":"8080"}
{"level":"info","time":"...","message":"Database repositories initialized"}
{"level":"info","time":"...","message":"Voice provider initialized","provider":"vapi"}
{"level":"info","time":"...","message":"Business services initialized"}
{"level":"info","time":"...","message":"HTTP router configured"}
{"level":"info","time":"...","message":"Server listening","port":"8080"}
```

## Step 5: Test the API

### Health Check
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### Register a Business
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "business_name": "Test Dental Clinic",
    "business_type": "dentist",
    "phone": "+1234567890",
    "email": "admin@testclinic.com",
    "password": "SecurePassword123!"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "user": {
    "id": "...",
    "business_id": "...",
    "email": "admin@testclinic.com",
    "role": "owner",
    "created_at": "..."
  }
}
```

Save the `access_token` for next requests!

### Get Business Details
```bash
export TOKEN="your_access_token_here"

curl -X GET http://localhost:8080/api/v1/businesses/me \
  -H "Authorization: Bearer $TOKEN"
```

### List Calls
```bash
curl -X GET http://localhost:8080/api/v1/calls \
  -H "Authorization: Bearer $TOKEN"
```

### Get Analytics
```bash
curl -X GET "http://localhost:8080/api/v1/analytics/overview?days=30" \
  -H "Authorization: Bearer $TOKEN"
```

## Step 6: Test with Vapi AI (Production)

Once you have a public URL (use ngrok for testing):

```bash
# Start ngrok
ngrok http 8080

# Update .env with ngrok URL
VAPI_WEBHOOK_URL=https://your-ngrok-url.ngrok.io/api/v1/webhooks/vapi
```

Now initiate a call:

```bash
curl -X POST http://localhost:8080/api/v1/calls \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+1234567890",
    "metadata": {
      "customer_name": "Jane Doe"
    }
  }'
```

---

## Using Docker

### Build and Run
```bash
docker-compose up --build
```

This starts:
- Application server on port 8080
- PostgreSQL database on port 5432

### Stop
```bash
docker-compose down
```

---

## Development Workflow

### Structure
```
.
├── cmd/server/           # Application entry point
├── internal/
│   ├── domain/          # Business entities & interfaces
│   ├── application/     # Services & business logic
│   ├── infrastructure/  # Database, providers
│   └── api/            # HTTP handlers
├── pkg/                 # Shared utilities
└── migrations/         # Database migrations
```

### Adding a New Feature

1. **Define domain entity** in `internal/domain/entities/`
2. **Create repository** in `internal/infrastructure/database/`
3. **Add service** in `internal/application/services/`
4. **Create handler** in `internal/api/handlers/`
5. **Update router** in `router.go`

### Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/application/services/
```

---

## Troubleshooting

### Database Connection Failed
- Verify `DATABASE_URL` in `.env`
- Check Supabase project is active
- Ensure migrations are run

### JWT Token Invalid
- Check `JWT_SECRET_KEY` is set
- Token may be expired (15 min for access tokens)
- Use refresh token endpoint to get new access token

### Vapi Webhook Not Working
- Verify `VAPI_WEBHOOK_URL` is publicly accessible
- Check webhook signature validation
- Look at server logs for webhook events

### Port Already in Use
- Change `SERVER_PORT` in `.env`
- Or stop the conflicting process:
  ```bash
  lsof -ti:8080 | xargs kill
  ```

---

## Next Steps

- Read the complete API documentation in `API.md`
- Review the architecture in `README.md`
- Check implementation progress in `PROGRESS.md`
- Write tests for your custom logic
- Deploy to production (Heroku, AWS, GCP, etc.)

---

## Need Help?

- Check logs: Application uses structured JSON logging
- Review error responses: All errors include error codes
- Enable debug logging: Set `LOG_LEVEL=debug` in `.env`

Happy coding! 🚀
