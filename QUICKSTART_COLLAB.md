# 🚀 Quick Start - Collaborative Notes

## Local Development Setup (5 minutes)

### 1. Database Setup
```bash
cd Backend
psql -U your_username -d projectnest -f migrations/003_add_collaborative_notes.sql
```

### 2. Start Collaboration Server
```bash
cd collab-server
npm install
npm start
```
✅ Server running at `ws://localhost:1234`

### 3. Configure Frontend
```bash
cd Frontend
echo "VITE_COLLAB_SERVER_URL=ws://localhost:1234" >> .env
npm install
npm run dev
```

### 4. Start Backend
```bash
cd Backend
./server
```

### 5. Test It Out
1. Open `http://localhost:5173` in two browser windows
2. Navigate to any project's notes
3. Create a new note
4. Start typing in one window
5. Watch it sync in real-time in the other window! ✨

## Features You Can Try

### Basic Formatting
- **Bold**: Ctrl/Cmd + B
- *Italic*: Ctrl/Cmd + I
- <u>Underline</u>: Ctrl/Cmd + U
- ~~Strikethrough~~: Ctrl/Cmd + Shift + S

### Headings
- Type `# ` + Space for H1
- Type `## ` + Space for H2
- Type `### ` + Space for H3

### Lists
- Type `- ` + Space for bullet list
- Type `1. ` + Space for numbered list
- Type `[ ] ` + Space for task list

### Code
- Inline code: Wrap text with backticks \`code\`
- Code block: Type \`\`\` + Space

### Links
- Click link icon or use Ctrl/Cmd + K

## Keyboard Shortcuts

| Action | Windows/Linux | Mac |
|--------|---------------|-----|
| Bold | Ctrl + B | Cmd + B |
| Italic | Ctrl + I | Cmd + I |
| Underline | Ctrl + U | Cmd + U |
| Strikethrough | Ctrl + Shift + S | Cmd + Shift + S |
| Code | Ctrl + E | Cmd + E |
| Undo | Ctrl + Z | Cmd + Z |
| Redo | Ctrl + Shift + Z | Cmd + Shift + Z |

## Troubleshooting

### ❌ Can't connect to collaboration server
**Solution**: Make sure collab server is running:
```bash
cd collab-server && npm start
```

### ❌ Changes not syncing
**Solution**: 
1. Check browser console for errors
2. Verify WebSocket connection (Network tab)
3. Restart collab server

### ❌ "Module not found" errors
**Solution**: 
```bash
cd Frontend && npm install
cd ../collab-server && npm install
```

## Next Steps

1. ✅ Deploy collaboration server (see COLLABORATIVE_NOTES_SETUP.md)
2. ✅ Configure authentication for WebSocket
3. ✅ Set up SSL/TLS for production (wss://)
4. ✅ Configure backups for yjs-storage

## Need Help?

- 📖 Full documentation: `COLLABORATIVE_NOTES_SETUP.md`
- 🐛 Issues: Check server logs
- 💬 Questions: Open GitHub issue

---

**That's it! You're ready to collaborate! 🎉**
