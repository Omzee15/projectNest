FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy application files
COPY . .

FROM node:20-alpine

WORKDIR /app

# Copy from builder
COPY --from=builder /app .

# Create storage directory
RUN mkdir -p /app/yjs-storage && \
    chown -R node:node /app

USER node

EXPOSE 5173

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:5173', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
