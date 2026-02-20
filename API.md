# CallPilotReceptionist API Documentation

Base URL: `http://localhost:8080`

## Authentication

All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

---

## Endpoints

### Health & Status

#### GET /health
Health check endpoint.

**Response**: 200 OK
```json
{"status": "ok"}
```

#### GET /ready
Readiness check endpoint.

**Response**: 200 OK
```json
{"status": "ready"}
```

---

### Authentication

#### POST /api/v1/auth/register
Register a new business and user.

**Request Body**:
```json
{
  "business_name": "John's Dental Clinic",
  "business_type": "dentist",
  "phone": "+1234567890",
  "email": "john@dentist.com",
  "password": "securepassword123"
}
```

**Response**: 201 Created
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "business_id": "uuid",
    "email": "john@dentist.com",
    "role": "owner",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### POST /api/v1/auth/login
Login with email and password.

**Request Body**:
```json
{
  "email": "john@dentist.com",
  "password": "securepassword123"
}
```

**Response**: 200 OK
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": { ... }
}
```

#### POST /api/v1/auth/refresh
Refresh access token using refresh token.

**Request Body**:
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response**: 200 OK
```json
{
  "access_token": "eyJhbGc..."
}
```

#### POST /api/v1/auth/logout
Logout (client should delete tokens).

**Headers**: `Authorization: Bearer <token>`

**Response**: 200 OK
```json
{
  "message": "Logged out successfully"
}
```

---

### Business Management

#### GET /api/v1/businesses/me
Get current business details.

**Headers**: `Authorization: Bearer <token>`

**Response**: 200 OK
```json
{
  "id": "uuid",
  "name": "John's Dental Clinic",
  "type": "dentist",
  "phone": "+1234567890",
  "settings": {},
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### PUT /api/v1/businesses/me
Update business information.

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "name": "John's Advanced Dental Clinic",
  "type": "dentist",
  "phone": "+1234567890",
  "settings": {
    "working_hours": "9AM-5PM"
  }
}
```

**Response**: 200 OK
```json
{
  "id": "uuid",
  "name": "John's Advanced Dental Clinic",
  ...
}
```

---

### Call Management

#### POST /api/v1/calls
Initiate a new call.

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "phone_number": "+1234567890",
  "assistant_id": "optional-assistant-id",
  "metadata": {
    "customer_name": "Jane Doe"
  }
}
```

**Response**: 201 Created
```json
{
  "id": "uuid",
  "business_id": "uuid",
  "provider_call_id": "vapi-call-id",
  "caller_phone": "+1234567890",
  "duration": 0,
  "status": "initiated",
  "cost": 0,
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### GET /api/v1/calls
List all calls for the business.

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `limit` (optional): Number of results (default: 20, max: 100)
- `offset` (optional): Pagination offset (default: 0)
- `status` (optional): Filter by status

**Response**: 200 OK
```json
{
  "calls": [
    {
      "id": "uuid",
      "business_id": "uuid",
      "caller_phone": "+1234567890",
      "duration": 120,
      "status": "completed",
      "cost": 0.05,
      "started_at": "2024-01-01T00:00:00Z",
      "ended_at": "2024-01-01T00:02:00Z",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

#### GET /api/v1/calls/:id
Get specific call details.

**Headers**: `Authorization: Bearer <token>`

**Response**: 200 OK
```json
{
  "id": "uuid",
  "business_id": "uuid",
  "provider_call_id": "vapi-call-id",
  "caller_phone": "+1234567890",
  "duration": 120,
  "status": "completed",
  "cost": 0.05,
  "started_at": "2024-01-01T00:00:00Z",
  "ended_at": "2024-01-01T00:02:00Z",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### GET /api/v1/calls/:id/transcript
Get call transcript.

**Headers**: `Authorization: Bearer <token>`

**Response**: 200 OK
```json
{
  "call_id": "uuid",
  "messages": [
    {
      "role": "assistant",
      "message": "Hello, how can I help you today?",
      "timestamp": "2024-01-01T00:00:05Z"
    },
    {
      "role": "user",
      "message": "I'd like to schedule an appointment",
      "timestamp": "2024-01-01T00:00:10Z"
    }
  ]
}
```

#### POST /api/v1/webhooks/vapi
Webhook endpoint for Vapi AI events (no auth required - signature validated).

**Headers**: `X-Vapi-Signature: <hmac-signature>`

**Request Body**: Raw JSON from Vapi
```json
{
  "type": "call.ended",
  "callId": "vapi-call-id",
  "status": "completed",
  ...
}
```

**Response**: 200 OK
```json
{
  "message": "Webhook processed successfully"
}
```

---

### Interactions & Appointments

#### GET /api/v1/calls/:id/interactions
Get interactions for a specific call.

**Headers**: `Authorization: Bearer <token>`

**Response**: 200 OK
```json
[
  {
    "id": "uuid",
    "call_id": "uuid",
    "type": "appointment_request",
    "content": {
      "date": "2024-01-15",
      "time": "10:00 AM",
      "service": "cleaning"
    },
    "timestamp": "2024-01-01T00:01:00Z",
    "created_at": "2024-01-01T00:01:05Z"
  }
]
```

#### GET /api/v1/interactions
List all interactions for the business.

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `limit` (optional): Number of results (default: 20, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response**: 200 OK
```json
{
  "interactions": [ ... ],
  "total": 50,
  "limit": 20,
  "offset": 0
}
```

#### GET /api/v1/appointments
List all appointments.

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `limit` (optional): Number of results (default: 20, max: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response**: 200 OK
```json
[
  {
    "id": "uuid",
    "call_id": "uuid",
    "business_id": "uuid",
    "customer_name": "Jane Doe",
    "customer_phone": "+1234567890",
    "requested_date": "2024-01-15",
    "requested_time": "10:00 AM",
    "service_type": "cleaning",
    "notes": "First time patient",
    "status": "pending",
    "extracted_at": "2024-01-01T00:01:00Z",
    "created_at": "2024-01-01T00:01:05Z"
  }
]
```

#### PATCH /api/v1/appointments/:id
Update appointment status.

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "status": "confirmed"
}
```

**Valid statuses**: `pending`, `confirmed`, `cancelled`, `completed`

**Response**: 200 OK
```json
{
  "id": "uuid",
  "status": "confirmed",
  "confirmed_at": "2024-01-01T00:05:00Z",
  ...
}
```

---

### Analytics

#### GET /api/v1/analytics/overview
Get dashboard overview statistics.

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `days` (optional): Number of days to analyze (default: 30)

**Response**: 200 OK
```json
{
  "total_calls": 150,
  "completed_calls": 140,
  "failed_calls": 10,
  "total_duration": 18000,
  "average_duration": 120.0,
  "total_cost": 7.50,
  "pending_appointments": 5
}
```

#### GET /api/v1/analytics/calls
Get call volume over time.

**Headers**: `Authorization: Bearer <token>`

**Query Parameters**:
- `days` (optional): Number of days (default: 30)

**Response**: 200 OK
```json
{
  "data": [
    {
      "date": "2024-01-01",
      "count": 5
    },
    {
      "date": "2024-01-02",
      "count": 8
    }
  ]
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message",
  "details": "Optional additional details"
}
```

**Common Error Codes**:
- `NOT_FOUND` - Resource not found (404)
- `ALREADY_EXISTS` - Resource already exists (409)
- `INVALID_INPUT` - Invalid request data (400)
- `VALIDATION_ERROR` - Validation failed (400)
- `UNAUTHORIZED` - Authentication required or invalid token (401)
- `FORBIDDEN` - Access denied (403)
- `PROVIDER_ERROR` - Voice provider error (502)
- `INTERNAL_ERROR` - Internal server error (500)

---

## Testing with cURL

### Register and Login
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "business_name": "Test Clinic",
    "business_type": "dentist",
    "phone": "+1234567890",
    "email": "test@clinic.com",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@clinic.com",
    "password": "password123"
  }'
```

### Use Protected Endpoints
```bash
# Get business details
curl -X GET http://localhost:8080/api/v1/businesses/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# List calls
curl -X GET "http://localhost:8080/api/v1/calls?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get analytics
curl -X GET "http://localhost:8080/api/v1/analytics/overview?days=7" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## Rate Limiting

Currently no rate limiting is implemented. This should be added in production.

## Pagination

List endpoints support pagination with `limit` and `offset` parameters:
- `limit`: Max 100, default 20
- `offset`: Starting position, default 0

---

## Next Steps

1. Set up your Supabase database
2. Run migrations from `migrations/001_initial_schema.up.sql`
3. Configure `.env` with your credentials
4. Start the server: `./server`
5. Test with the examples above!
