#!/bin/bash

# ProjectNest Collaborative Notes - Automated Setup Script
# This script automates the installation of dependencies for the collaborative notes system

set -e  # Exit on error

echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║   ProjectNest Collaborative Notes - Automated Setup          ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Check prerequisites
echo "Checking prerequisites..."

# Check Node.js
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js 18+ first."
    exit 1
fi
NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
    print_error "Node.js version must be 18 or higher. Current: $(node -v)"
    exit 1
fi
print_success "Node.js $(node -v) detected"

# Check npm
if ! command -v npm &> /dev/null; then
    print_error "npm is not installed. Please install npm first."
    exit 1
fi
print_success "npm $(npm -v) detected"

# Check PostgreSQL
if ! command -v psql &> /dev/null; then
    print_info "PostgreSQL CLI not detected. Make sure PostgreSQL is installed and accessible."
else
    print_success "PostgreSQL detected"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "Step 1: Installing Collaboration Server Dependencies"
echo "═══════════════════════════════════════════════════════════════"

if [ -d "$SCRIPT_DIR/collab-server" ]; then
    cd "$SCRIPT_DIR/collab-server"
    print_info "Installing collaboration server dependencies..."
    npm install
    print_success "Collaboration server dependencies installed"
    
    # Create .env if it doesn't exist
    if [ ! -f ".env" ]; then
        print_info "Creating .env file from template..."
        cp .env.example .env
        print_success ".env file created"
    else
        print_info ".env file already exists"
    fi
    
    # Create yjs-storage directory
    if [ ! -d "yjs-storage" ]; then
        mkdir -p yjs-storage
        print_success "Created yjs-storage directory"
    fi
else
    print_error "collab-server directory not found!"
    exit 1
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "Step 2: Installing Frontend Dependencies"
echo "═══════════════════════════════════════════════════════════════"

if [ -d "$SCRIPT_DIR/Frontend" ]; then
    cd "$SCRIPT_DIR/Frontend"
    print_info "Installing frontend dependencies (this may take a few minutes)..."
    npm install
    print_success "Frontend dependencies installed"
    
    # Update .env
    if [ -f ".env" ]; then
        if ! grep -q "VITE_COLLAB_SERVER_URL" .env; then
            echo "" >> .env
            echo "# Collaboration Server" >> .env
            echo "VITE_COLLAB_SERVER_URL=ws://localhost:1234" >> .env
            print_success "Added VITE_COLLAB_SERVER_URL to .env"
        else
            print_info "VITE_COLLAB_SERVER_URL already exists in .env"
        fi
    else
        print_info "Creating .env file..."
        cp .env.example .env
        echo "" >> .env
        echo "# Collaboration Server" >> .env
        echo "VITE_COLLAB_SERVER_URL=ws://localhost:1234" >> .env
        print_success ".env file created"
    fi
else
    print_error "Frontend directory not found!"
    exit 1
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "Step 3: Database Migration"
echo "═══════════════════════════════════════════════════════════════"

if [ -f "$SCRIPT_DIR/Backend/migrations/003_add_collaborative_notes.sql" ]; then
    print_info "Database migration file found"
    echo ""
    read -p "Do you want to run the database migration now? (y/n) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter PostgreSQL database name: " DB_NAME
        read -p "Enter PostgreSQL username: " DB_USER
        
        cd "$SCRIPT_DIR/Backend"
        if psql -U "$DB_USER" -d "$DB_NAME" -f migrations/003_add_collaborative_notes.sql; then
            print_success "Database migration completed successfully"
        else
            print_error "Database migration failed. Please run it manually."
        fi
    else
        print_info "Skipping database migration. Run it manually later:"
        echo "  cd Backend"
        echo "  psql -U your_user -d your_db -f migrations/003_add_collaborative_notes.sql"
    fi
else
    print_error "Migration file not found!"
fi

echo ""
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║                    Installation Complete! 🎉                  ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""
echo "Next steps:"
echo ""
echo "1. Start the collaboration server:"
echo "   ${GREEN}cd collab-server && npm start${NC}"
echo ""
echo "2. Start the frontend (in another terminal):"
echo "   ${GREEN}cd Frontend && npm run dev${NC}"
echo ""
echo "3. Start the backend (in another terminal):"
echo "   ${GREEN}cd Backend && ./server${NC}"
echo ""
echo "4. Open http://localhost:5173 and test collaborative editing!"
echo ""
echo "📚 Documentation:"
echo "   - Quick Start: QUICKSTART_COLLAB.md"
echo "   - Full Setup: COLLABORATIVE_NOTES_SETUP.md"
echo "   - Summary: IMPLEMENTATION_SUMMARY.md"
echo ""
echo "🐛 Troubleshooting:"
echo "   - Check server logs if WebSocket connection fails"
echo "   - Verify all services are running on correct ports"
echo "   - See COLLABORATIVE_NOTES_SETUP.md for common issues"
echo ""
print_success "Setup script completed successfully!"
