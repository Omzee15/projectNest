# 🚀 ProjectNest Collaborative Notes System

> **Google Docs-style collaborative editing for ProjectNest with rich text formatting**

## ✨ What's New

ProjectNest notes now support **real-time collaborative editing** with rich text features, powered by CRDT technology for conflict-free collaboration.

### Key Features

🎨 **Rich Text Editing**
- Bold, Italic, Underline, Strikethrough
- Headings (H1-H6)
- Bullet & Numbered Lists
- Task Lists with checkboxes
- Code blocks & inline code
- Links & Blockquotes

👥 **Real-time Collaboration**
- Multi-user editing
- Live cursor tracking
- User presence indicators
- Automatic conflict resolution
- Offline support with auto-sync

💾 **Smart Persistence**
- Auto-save (every 3 seconds)
- Version tracking
- Document snapshots
- LevelDB persistence

## 📦 What's Included

```
projectNest/
├── Backend/
│   ├── migrations/
│   │   └── 003_add_collaborative_notes.sql    # Database migration
│   └── internal/
│       └── models/                             # Updated Go models
├── collab-server/                              # ⭐ New WebSocket Server
│   ├── server.js                               # Main server
│   ├── package.json                            # Dependencies
│   ├── Dockerfile                              # Container config
│   └── README.md                               # Server docs
├── Frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── CollaborativeEditor.tsx         # ⭐ Rich text editor
│   │   │   ├── CollaboratorPresence.tsx        # ⭐ User presence
│   │   │   └── ProjectNotesCollaborative.tsx   # ⭐ Enhanced notes UI
│   │   ├── styles/
│   │   │   └── editor.css                      # ⭐ Editor styling
│   │   └── types/
│   │       └── index.ts                        # Updated types
│   └── package.json                            # New dependencies
├── COLLABORATIVE_NOTES_SETUP.md                # 📖 Full documentation
├── QUICKSTART_COLLAB.md                        # 🚀 Quick start guide
├── IMPLEMENTATION_SUMMARY.md                   # 📊 What was built
└── setup-collaborative-notes.sh                # 🛠️ Automated setup
```

## ⚡ Quick Start

### Automated Setup (Recommended)

```bash
# Run the automated setup script
./setup-collaborative-notes.sh
```

The script will:
1. ✅ Install all dependencies
2. ✅ Configure environment files
3. ✅ Create necessary directories
4. ✅ Optionally run database migration

### Manual Setup

#### 1. Install Dependencies

```bash
# Collaboration server
cd collab-server
npm install

# Frontend
cd ../Frontend
npm install
```

#### 2. Database Migration

```bash
cd Backend
psql -U your_user -d projectnest -f migrations/003_add_collaborative_notes.sql
```

#### 3. Configure Environment

**Collaboration Server** (`collab-server/.env`):
```env
COLLAB_PORT=1234
COLLAB_HOST=0.0.0.0
PERSISTENCE_DIR=./yjs-storage
```

**Frontend** (`Frontend/.env`):
```env
VITE_COLLAB_SERVER_URL=ws://localhost:1234
```

#### 4. Start Services

```bash
# Terminal 1 - Collaboration Server
cd collab-server
npm start

# Terminal 2 - Frontend
cd Frontend
npm run dev

# Terminal 3 - Backend
cd Backend
./server
```

#### 5. Test It!

1. Open `http://localhost:5173` in **two browser windows**
2. Navigate to any project's notes
3. Create a new note
4. Start typing in one window
5. Watch it sync in real-time in the other! ✨

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [QUICKSTART_COLLAB.md](QUICKSTART_COLLAB.md) | Get started in 5 minutes |
| [COLLABORATIVE_NOTES_SETUP.md](COLLABORATIVE_NOTES_SETUP.md) | Complete setup & deployment guide |
| [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) | Technical details & architecture |
| [collab-server/README.md](collab-server/README.md) | WebSocket server documentation |

## 🏗️ Architecture

```
┌──────────────┐
│   Browser 1  │───┐
└──────────────┘   │
                   │  WebSocket
┌──────────────┐   │  (Real-time)
│   Browser 2  │───┼──────────────► ┌─────────────────┐
└──────────────┘   │                │  Collab Server  │
                   │                │  (Yjs + Node)   │
┌──────────────┐   │                └────────┬────────┘
│   Browser 3  │───┘                         │
└──────────────┘                             │
                                             ▼
      │                              ┌──────────────┐
      │ HTTP/REST API                │   LevelDB    │
      ▼                              │ (Yjs Storage)│
┌─────────────┐                      └──────────────┘
│ Go Backend  │
│  (REST API) │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ PostgreSQL  │
│  (Metadata) │
└─────────────┘
```

### Why This Architecture?

- **Separation of Concerns**: Real-time sync is handled by specialized WebSocket server
- **Scalability**: Can scale collab server independently
- **Reliability**: If collab server goes down, backend API still works
- **Performance**: CRDT in browser, minimal server processing
- **Flexibility**: Easy to deploy on different infrastructure

## 🎯 Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Editor | TipTap 2.x | Rich text editing |
| CRDT | Yjs | Conflict-free replication |
| Transport | y-websocket | Real-time sync |
| Server | Node.js + Express | WebSocket server |
| Storage | LevelDB | Document persistence |
| Backend | Go + Gin | REST API |
| Database | PostgreSQL | Metadata storage |
| Frontend | React + TypeScript | UI framework |

## 🔧 Configuration

### Auto-Save Interval

Default: 3 seconds (3000ms)

Adjust in `ProjectNotesCollaborative.tsx`:
```typescript
<CollaborativeEditor
  autoSaveInterval={5000}  // 5 seconds
  // ...
/>
```

### WebSocket URL

Development: `ws://localhost:1234`
Production: `wss://your-server.com`

Set in `Frontend/.env`:
```env
VITE_COLLAB_SERVER_URL=ws://localhost:1234
```

### Connection Timeout

Default: 60 seconds

Adjust in `collab-server/server.js`:
```javascript
const timeout = 60000; // milliseconds
```

## 🚢 Deployment

### Docker (Recommended)

```bash
cd collab-server
docker build -t projectnest-collab .
docker run -d -p 1234:1234 \
  -v $(pwd)/yjs-storage:/app/yjs-storage \
  projectnest-collab
```

### Render.com / Railway

1. Create new Web Service
2. Build: `cd collab-server && npm install`
3. Start: `cd collab-server && npm start`
4. Add environment variables
5. Update frontend `.env` with deployed URL

See [COLLABORATIVE_NOTES_SETUP.md](COLLABORATIVE_NOTES_SETUP.md) for detailed deployment instructions.

## 🔒 Security

### Production Checklist

- [ ] Enable WebSocket authentication (JWT tokens)
- [ ] Configure CORS properly
- [ ] Use WSS (secure WebSocket) in production
- [ ] Set up rate limiting
- [ ] Configure firewall rules
- [ ] Enable HTTPS for all services
- [ ] Set up monitoring and alerts
- [ ] Configure automated backups

See security section in [COLLABORATIVE_NOTES_SETUP.md](COLLABORATIVE_NOTES_SETUP.md)

## 📊 Performance

### Expected Performance

- **Concurrent users per document**: 50+
- **Document size**: Up to 1MB comfortably
- **Sync latency**: <100ms (good connection)
- **Memory per doc**: ~50MB
- **Storage per doc**: ~2KB (LevelDB)

### Optimization Tips

1. **Adjust auto-save interval** for large documents
2. **Enable WebSocket compression** in production
3. **Configure LevelDB cache size**
4. **Use CDN** for static assets
5. **Enable HTTP/2** on server

## 🐛 Troubleshooting

### Can't connect to WebSocket

```bash
# Check if server is running
curl http://localhost:1234/health

# Check firewall
lsof -i :1234
```

### Changes not syncing

1. Open browser console
2. Check for WebSocket errors
3. Verify document ID format
4. Restart collab server

### High memory usage

- Reduce document timeout
- Increase cleanup frequency
- Monitor active documents: `curl http://localhost:1234/docs`

See full troubleshooting guide in [COLLABORATIVE_NOTES_SETUP.md](COLLABORATIVE_NOTES_SETUP.md)

## 📈 Monitoring

### Health Check

```bash
curl http://localhost:1234/health
```

Response:
```json
{
  "status": "ok",
  "activeConnections": 5,
  "activeDocs": 12
}
```

### Active Documents

```bash
curl http://localhost:1234/docs
```

### Document Info

```bash
curl http://localhost:1234/docs/project-123-note-456
```

## 🎓 Learning Resources

- [TipTap Documentation](https://tiptap.dev/)
- [Yjs Documentation](https://docs.yjs.dev/)
- [CRDT Technology](https://crdt.tech/)
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

## 🤝 Contributing

Contributions are welcome! Areas for improvement:

- [ ] Add more TipTap extensions (tables, images, etc.)
- [ ] Implement version history UI
- [ ] Add export functionality (PDF, Markdown)
- [ ] Enhance mobile experience
- [ ] Add keyboard shortcuts customization
- [ ] Implement @mentions
- [ ] Add collaborative comments

## 📄 License

Same as ProjectNest main project.

---

## 🎉 What's Next?

1. **Try it locally** - Run the setup script and test it out
2. **Deploy to production** - Follow the deployment guide
3. **Customize** - Add your own TipTap extensions
4. **Contribute** - Share improvements with the community

**Need help?** Check the documentation or open an issue on GitHub.

---

**Built with ❤️ for better collaboration**
