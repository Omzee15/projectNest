# Collaboration Server

WebSocket server for real-time collaborative editing in ProjectNest using Yjs CRDT.

## Features

- Real-time collaborative editing with conflict-free synchronization
- Persistent document storage using LevelDB
- User awareness and cursor tracking
- Automatic document cleanup
- Health monitoring endpoints

## Installation

```bash
npm install
```

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

## Running the Server

### Development
```bash
npm run dev
```

### Production
```bash
npm start
```

## Endpoints

### WebSocket
- `ws://localhost:1234/{documentId}` - Connect to a document

### HTTP API
- `GET /health` - Health check
- `GET /docs` - List all active documents
- `GET /docs/:docName` - Get document info
- `POST /snapshot/:docName` - Create document snapshot

## Document Naming Convention

Documents are identified by: `project-{projectUid}-note-{noteUid}`

Example: `project-123e4567-e89b-note-987f6543-a21b`

## Deployment

### Docker (Recommended)

```bash
docker build -t projectnest-collab .
docker run -p 1234:1234 -v $(pwd)/yjs-storage:/app/yjs-storage projectnest-collab
```

### Render.com / Railway

1. Connect your repository
2. Set build command: `cd collab-server && npm install`
3. Set start command: `cd collab-server && npm start`
4. Add environment variables from `.env.example`

## Monitoring

Check server health:
```bash
curl http://localhost:1234/health
```

View active documents:
```bash
curl http://localhost:1234/docs
```

## Security Considerations

1. Add authentication middleware in production
2. Configure CORS properly
3. Use secure WebSocket (wss://) in production
4. Implement rate limiting
5. Set up proper firewall rules
