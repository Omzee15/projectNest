# Collaborative Notes System - Setup & Deployment Guide

## 🎯 Overview

ProjectNest now features a real-time collaborative notes system powered by:
- **TipTap** - Rich text editor with extensive formatting
- **Yjs** - CRDT (Conflict-free Replicated Data Type) for seamless collaboration
- **y-websocket** - Real-time synchronization
- **LevelDB** - Persistent document storage

## 📋 Features

### Rich Text Editing
- ✅ Bold, Italic, Underline, Strikethrough
- ✅ Headings (H1-H6)
- ✅ Bullet lists, Numbered lists
- ✅ Task lists with checkboxes
- ✅ Links and code blocks
- ✅ Blockquotes
- ✅ Undo/Redo

### Collaborative Features
- ✅ Real-time multi-user editing
- ✅ User presence indicators
- ✅ Cursor awareness (see where others are typing)
- ✅ Automatic conflict resolution (CRDT)
- ✅ Auto-save (configurable interval)
- ✅ Offline support with automatic sync
- ✅ Version tracking

## 🚀 Quick Start

### 1. Database Migration

Run the migration to add collaborative editing support:

```bash
cd Backend
psql -U your_username -d your_database -f migrations/003_add_collaborative_notes.sql
```

Or apply through your database management tool.

### 2. Install Dependencies

#### Collaboration Server
```bash
cd collab-server
npm install
```

#### Frontend
```bash
cd Frontend
npm install
```

### 3. Configure Environment Variables

#### Collaboration Server (.env)
```bash
cd collab-server
cp .env.example .env
```

Edit `.env`:
```env
COLLAB_PORT=1234
COLLAB_HOST=0.0.0.0
PERSISTENCE_DIR=./yjs-storage
NODE_ENV=production
```

#### Frontend (.env)
```bash
cd Frontend
cp .env.example .env
```

Edit `.env`:
```env
# For local development
VITE_COLLAB_SERVER_URL=ws://localhost:1234

# For production
# VITE_COLLAB_SERVER_URL=wss://your-collab-server.com
```

### 4. Start the Servers

#### Start Collaboration Server
```bash
cd collab-server
npm start
```

The server will start on `ws://localhost:1234`

#### Start Frontend (in another terminal)
```bash
cd Frontend
npm run dev
```

#### Start Backend (in another terminal)
```bash
cd Backend
./server  # or go run cmd/server/main.go
```

## 🏗️ Architecture

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│   Browser   │◄───────►│  Collab Server   │◄───────►│   LevelDB   │
│  (TipTap +  │  WebSocket │  (y-websocket)  │         │ (Persistent │
│    Yjs)     │         │                  │         │  Storage)   │
└─────────────┘         └──────────────────┘         └─────────────┘
      │
      │ REST API
      │
      ▼
┌─────────────┐         ┌──────────────────┐
│  Go Backend │◄───────►│   PostgreSQL     │
│  (REST API) │         │   (Metadata)     │
└─────────────┘         └──────────────────┘
```

### Data Flow

1. **User opens note** → Frontend creates Yjs document and connects to WebSocket
2. **User types** → Changes sync via WebSocket to all connected users
3. **Auto-save** → HTML content saved to PostgreSQL via REST API
4. **WebSocket server** → Persists Yjs state to LevelDB for recovery
5. **User reconnects** → Document state restored from LevelDB

## 🔧 Configuration

### Auto-Save Settings

In `CollaborativeEditor.tsx`:
```typescript
<CollaborativeEditor
  autoSave={true}
  autoSaveInterval={3000}  // milliseconds (3 seconds)
  // ...
/>
```

### WebSocket Connection

Document IDs follow the pattern:
```
project-{projectUid}-note-{noteUid}
```

Example: `project-123e4567-e89b-note-987f6543-a21b`

## 🚢 Production Deployment

### Option 1: Docker Deployment

#### Build and Run Collaboration Server
```bash
cd collab-server
docker build -t projectnest-collab .
docker run -d \
  -p 1234:1234 \
  -v $(pwd)/yjs-storage:/app/yjs-storage \
  --name projectnest-collab \
  projectnest-collab
```

### Option 2: Render.com / Railway

#### Collaboration Server

1. **Create New Web Service**
2. **Connect Repository**
3. **Settings:**
   - Build Command: `cd collab-server && npm install`
   - Start Command: `cd collab-server && npm start`
   - Port: `1234`

4. **Environment Variables:**
   ```
   COLLAB_PORT=1234
   COLLAB_HOST=0.0.0.0
   PERSISTENCE_DIR=./yjs-storage
   NODE_ENV=production
   ```

5. **Add Persistent Disk** (for yjs-storage)

6. **Update Frontend .env:**
   ```
   VITE_COLLAB_SERVER_URL=wss://your-service.onrender.com
   ```

### Option 3: Heroku

```bash
cd collab-server
heroku create projectnest-collab
heroku config:set COLLAB_PORT=1234
git push heroku main
```

## 🔒 Security Considerations

### 1. Authentication

Add authentication to WebSocket connections in `server.js`:

```javascript
wss.on('connection', (conn, req) => {
  // Extract token from query params
  const url = new URL(req.url, `http://${req.headers.host}`);
  const token = url.searchParams.get('token');
  
  // Verify token
  if (!verifyToken(token)) {
    conn.close(1008, 'Unauthorized');
    return;
  }
  
  // Continue with setup...
});
```

### 2. CORS Configuration

Update `server.js`:
```javascript
const allowedOrigins = process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:5173'];

app.use(cors({
  origin: (origin, callback) => {
    if (!origin || allowedOrigins.includes(origin)) {
      callback(null, true);
    } else {
      callback(new Error('Not allowed by CORS'));
    }
  }
}));
```

### 3. Rate Limiting

Install and configure:
```bash
npm install express-rate-limit
```

```javascript
const rateLimit = require('express-rate-limit');

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100 // limit each IP to 100 requests per windowMs
});

app.use(limiter);
```

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:1234/health
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2024-11-20T10:30:00.000Z",
  "activeConnections": 5,
  "activeDocs": 12
}
```

### Active Documents
```bash
curl http://localhost:1234/docs
```

### Specific Document Info
```bash
curl http://localhost:1234/docs/project-123-note-456
```

## 🐛 Troubleshooting

### WebSocket Connection Fails

1. **Check server is running:**
   ```bash
   curl http://localhost:1234/health
   ```

2. **Verify port is open:**
   ```bash
   lsof -i :1234
   ```

3. **Check browser console** for WebSocket errors

4. **Firewall rules** - Ensure port 1234 is open

### Document Not Syncing

1. **Check browser console** for errors
2. **Verify document ID** matches pattern
3. **Check network tab** for WebSocket frames
4. **Restart collaboration server**

### High Memory Usage

1. **Increase cleanup frequency** in `server.js`:
   ```javascript
   setInterval(() => {
     // cleanup code
   }, 15000); // Reduce from 30s to 15s
   ```

2. **Reduce document timeout**:
   ```javascript
   const timeout = 30000; // Reduce from 60s to 30s
   ```

## 🔄 Backup & Recovery

### Backup Yjs Storage
```bash
cd collab-server
tar -czf yjs-storage-backup-$(date +%Y%m%d).tar.gz yjs-storage/
```

### Restore from Backup
```bash
cd collab-server
tar -xzf yjs-storage-backup-20241120.tar.gz
npm start
```

### Database Backup (PostgreSQL)
```bash
pg_dump -U username -d projectnest > backup-$(date +%Y%m%d).sql
```

## 📈 Performance Tuning

### 1. WebSocket Compression
Enable in `server.js`:
```javascript
const wss = new WebSocketServer({ 
  server,
  perMessageDeflate: true
});
```

### 2. LevelDB Settings
```javascript
const ldb = new leveldb.LeveldbPersistence(PERSISTENCE_DIR, {
  cacheSize: 16 * 1024 * 1024, // 16MB cache
  compression: true
});
```

### 3. Auto-Save Throttling
Increase interval for large documents:
```typescript
autoSaveInterval={5000} // 5 seconds instead of 3
```

## 📝 Usage Examples

### Basic Note Editing
```typescript
import { ProjectNotesCollaborative } from '@/components/ProjectNotesCollaborative';

function NotesPage() {
  return (
    <ProjectNotesCollaborative
      projectUid={projectUid}
      notes={notes}
      onCreateNote={handleCreate}
      onUpdateNote={handleUpdate}
      onDeleteNote={handleDelete}
    />
  );
}
```

### Read-Only Mode
```typescript
<ProjectNotesCollaborative
  projectUid={projectUid}
  notes={notes}
  readOnly={true}
/>
```

### Custom Auto-Save
```typescript
<CollaborativeEditor
  documentId={docId}
  userName={user.name}
  onSave={(content, html) => {
    // Custom save logic
    api.saveNote({ content, html });
  }}
  autoSave={true}
  autoSaveInterval={10000} // 10 seconds
/>
```

## 🎨 Customization

### Custom User Colors
Modify `getUserColor()` in `ProjectNotesCollaborative.tsx`:
```typescript
const getUserColor = (userId: string): string => {
  const colors = ['#3B82F6', '#10B981', /* add more colors */];
  const hash = userId.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
  return colors[hash % colors.length];
};
```

### Custom Toolbar
Modify `MenuBar` component in `CollaborativeEditor.tsx` to add/remove buttons.

### Custom Placeholders
```typescript
<CollaborativeEditor
  placeholder="Start brainstorming ideas..."
  // ...
/>
```

## 📚 Additional Resources

- [TipTap Documentation](https://tiptap.dev/)
- [Yjs Documentation](https://docs.yjs.dev/)
- [CRDT Explained](https://crdt.tech/)
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

## 🤝 Support

For issues or questions:
1. Check the troubleshooting section
2. Review server logs: `docker logs projectnest-collab`
3. Open an issue on GitHub
4. Contact support team

## 📜 License

Same as ProjectNest main project license.
