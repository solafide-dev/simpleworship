#!/bin/bash

# This script is used to build the display server frontend
# This is hacked into the wails prebuild hooks to make sure the frontend is built before the server is built

echo "Building display frontend..."
echo "Durectory: $(dirname "$0")"
cd "$(dirname "$0")"
echo "Installing NPM Package"
npm install
echo "Building frontend"
npm run build