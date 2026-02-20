# Security Notes

## Secrets Management

All sensitive credentials are stored in `.env` file which is:
- ✅ Listed in `.gitignore`
- ✅ Never committed to the repository
- ✅ Only exists locally

### Protected Secrets:
- Database connection string (Supabase)
- JWT secret key
- Vapi API key
- Cloudflare tunnel URLs (temporary)

### Safe to Commit:
- `.env.example` - Template with placeholder values
- Source code - No hardcoded secrets
- Configuration structure - Only references to environment variables

## Setup for New Developers

1. Copy `.env.example` to `.env`
2. Fill in your own credentials:
   - Get Supabase connection string from your project
   - Generate JWT secret: `openssl rand -base64 32`
   - Get Vapi API key from dashboard.vapi.ai
   - Set up Cloudflare tunnel for local development

## Never Commit:
- `.env` files
- Database passwords
- API keys
- JWT secrets
- Private keys
