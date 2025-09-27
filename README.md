# # Lucid Lists Backend

A Go-based REST API backend for the Lucid Lists project management application.

## Architecture

The backend follows a clean architecture pattern with the following layers:

- **Handlers**: HTTP request/response handling, input validation, and JSON serialization
- **Services**: Business logic and orchestration between repositories
- **Repositories**: Database operations and data persistence
- **Models**: Data structures for database entities and API DTOs

## Features

- ✅ Project management (CRUD operations)
- ✅ List management within projects (CRUD + positioning)
- ✅ Task management within lists (CRUD + moving between lists)
- ✅ UUID-based external identifiers (internal IDs hidden from frontend)
- ✅ Soft deletes using `is_active` flags
- ✅ Structured logging with logrus
- ✅ CORS support for frontend integration
- ✅ Input validation with detailed error messages
- ✅ Graceful shutdown handling

## API Endpoints

### Projects
- `GET /api/projects` - List all active projects
- `GET /api/projects/{project_uid}` - Get project with lists and tasks
- `POST /api/projects` - Create new project
- `PUT /api/projects/{project_uid}` - Update project
- `DELETE /api/projects/{project_uid}` - Soft delete project

### Lists
- `POST /api/lists` - Create list in project
- `PUT /api/lists/{list_uid}` - Update list name
- `DELETE /api/lists/{list_uid}` - Delete list
- `PUT /api/lists/{list_uid}/position` - Update list position

### Tasks
- `POST /api/tasks` - Create task in list
- `PUT /api/tasks/{task_uid}` - Update task
- `DELETE /api/tasks/{task_uid}` - Delete task
- `POST /api/tasks/{task_uid}/move` - Move task to different list

### Health Check
- `GET /health` - Health check endpoint

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (optional, for easy PostgreSQL setup)

### Database Setup

1. **Using Docker Compose (Recommended)**:
```bash
cd lucid-lists-backend-
docker-compose up -d
```

2. **Manual PostgreSQL Setup**:
- Create a database named `lucid_lists`
- Run the SQL schema from `db_schema.sql`

### Environment Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Update the `.env` file with your database credentials and preferences:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=lucid_lists
DB_SSLMODE=disable

SERVER_PORT=8080
SERVER_HOST=localhost

APP_ENV=development
LOG_LEVEL=info

CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

### Running the Application

1. **Install dependencies**:
```bash
go mod tidy
```

2. **Run the server**:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` by default.

## API Response Format

All API responses follow this consistent structure:

**Success Response:**
```json
{
  "data": {...},
  "success": true,
  "message": "optional success message"
}
```

**Error Response:**
```json
{
  "error": "error_type",
  "message": "Human readable error message",
  "status_code": 400
}
```

## Example API Usage

### Create a Project
```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My New Project",
    "description": "A project for testing",
    "status": "active"
  }'
```

### Create a List in Project
```bash
curl -X POST http://localhost:8080/api/lists \
  -H "Content-Type: application/json" \
  -d '{
    "project_uid": "123e4567-e89b-12d3-a456-426614174000",
    "name": "To Do",
    "position": 1
  }'
```

### Create a Task in List
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "list_uid": "123e4567-e89b-12d3-a456-426614174001",
    "title": "Implement user authentication",
    "description": "Add login and registration functionality",
    "priority": "high",
    "status": "todo"
  }'
```

### Move Task to Different List
```bash
curl -X POST http://localhost:8080/api/tasks/123e4567-e89b-12d3-a456-426614174010/move \
  -H "Content-Type: application/json" \
  -d '{
    "list_uid": "123e4567-e89b-12d3-a456-426614174002"
  }'
```

## Frontend Integration

The backend is designed to work seamlessly with the existing frontend:

1. **All entities use UUID-based external identifiers**: `project_uid`, `list_uid`, `task_uid`
2. **Internal database IDs are never exposed** to the frontend
3. **API response structure matches** the frontend's expected format
4. **CORS is configured** for the development frontend (Vite dev server)

## Development

### Project Structure
```
lucid-lists-backend-/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models and DTOs
│   ├── repositories/    # Data access layer
│   ├── services/        # Business logic layer
│   └── utils/           # Utilities and helpers
├── pkg/logger/          # Logging configuration
├── db_schema.sql        # Database schema
├── docker-compose.yml   # PostgreSQL setup
└── .env.example        # Environment template
```

### Database Schema

The application uses the following main entities:

- **Projects**: Top-level containers for organizing work
- **Lists**: Columns within projects (like "To Do", "In Progress", "Done")
- **Tasks**: Individual work items within lists

All entities support soft deletes via the `is_active` field and include audit timestamps.

### Logging

The application uses structured logging with logrus. Log levels can be configured via the `LOG_LEVEL` environment variable:
- `debug`: Detailed debugging information
- `info`: General operational messages (default)
- `warn`: Warning messages
- `error`: Error conditions

## Next Steps

1. **Testing**: Add comprehensive unit and integration tests
2. **Authentication**: Implement user authentication and authorization
3. **Validation**: Add more sophisticated business rule validation
4. **Performance**: Add database query optimization and caching
5. **Documentation**: Generate OpenAPI/Swagger documentation
6. **Deployment**: Add Docker containerization and deployment configs

## Contributing

1. Follow Go best practices and conventions
2. Add appropriate logging for debugging and monitoring
3. Validate all inputs and handle errors gracefully
4. Write tests for new functionality
5. Update documentation for API changes
backend service for the lucid list project manager
