import dotenv from 'dotenv';
import http from 'http';
import express from 'express';
import cors from 'cors';
import { WebSocketServer } from 'ws';
import * as Y from 'yjs';
import * as syncProtocol from 'y-protocols/sync';
import * as awarenessProtocol from 'y-protocols/awareness';
import * as encoding from 'lib0/encoding';
import * as decoding from 'lib0/decoding';
import * as map from 'lib0/map';

dotenv.config();

const PORT = process.env.COLLAB_PORT || 1234;
const HOST = process.env.COLLAB_HOST || '0.0.0.0';

const wsReadyStateConnecting = 0;
const wsReadyStateOpen = 1;

const messageSync = 0;
const messageAwareness = 1;

const docs = new Map();

const send = (conn, message) => {
  if (conn.readyState !== wsReadyStateConnecting && conn.readyState !== wsReadyStateOpen) {
    conn.close();
  }
  try {
    conn.send(message, err => { err != null && conn.close(); });
  } catch (e) {
    conn.close();
  }
};

const setupWSConnection = (conn, req, docName) => {
  conn.binaryType = 'arraybuffer';
  
  const doc = map.setIfUndefined(docs, docName, () => new Y.Doc());
  doc.conns = doc.conns || new Set();
  doc.conns.add(conn);
  
  const awareness = doc.awareness || new awarenessProtocol.Awareness(doc);
  doc.awareness = awareness;
  
  const awarenessHandler = ({ added, updated, removed }) => {
    const changedClients = added.concat(updated, removed);
    if (conn.readyState === wsReadyStateOpen) {
      const encoder = encoding.createEncoder();
      encoding.writeVarUint(encoder, messageAwareness);
      encoding.writeVarUint8Array(encoder, awarenessProtocol.encodeAwarenessUpdate(awareness, changedClients));
      send(conn, encoding.toUint8Array(encoder));
    }
  };
  
  awareness.on('update', awarenessHandler);
  
  conn.on('message', message => {
    try {
      const encoder = encoding.createEncoder();
      const decoder = decoding.createDecoder(new Uint8Array(message));
      const messageType = decoding.readVarUint(decoder);
      
      if (messageType === messageSync) {
        encoding.writeVarUint(encoder, messageSync);
        syncProtocol.readSyncMessage(decoder, encoder, doc, conn);
        if (encoding.length(encoder) > 1) send(conn, encoding.toUint8Array(encoder));
      } else if (messageType === messageAwareness) {
        awarenessProtocol.applyAwarenessUpdate(awareness, decoding.readVarUint8Array(decoder), conn);
      }
    } catch (err) {
      console.error('Message error:', err);
    }
  });
  
  const encoder = encoding.createEncoder();
  encoding.writeVarUint(encoder, messageSync);
  syncProtocol.writeSyncStep1(encoder, doc);
  send(conn, encoding.toUint8Array(encoder));
  
  conn.on('close', () => {
    doc.conns.delete(conn);
    awareness.off('update', awarenessHandler);
    awarenessProtocol.removeAwarenessStates(awareness, [conn], null);
  });
};

const app = express();
app.use(cors());
app.use(express.json());

app.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    timestamp: new Date().toISOString(),
    activeDocs: docs.size,
  });
});

const server = http.createServer(app);
const wss = new WebSocketServer({ server });

wss.on('connection', (conn, req) => {
  const docName = req.url.slice(1);
  console.log(`[${new Date().toISOString()}] Connection to: ${docName}`);
  setupWSConnection(conn, req, docName);
});

server.listen(PORT, HOST, () => {
  console.log('='.repeat(60));
  console.log('🚀 ProjectNest Collaboration Server');
  console.log('='.repeat(60));
  console.log(`📡 WebSocket: ws://${HOST}:${PORT}`);
  console.log(`🏥 Health: http://${HOST}:${PORT}/health`);
  console.log('='.repeat(60));
});

process.on('SIGINT', () => {
  console.log('\n🛑 Shutting down...');
  process.exit(0);
});
