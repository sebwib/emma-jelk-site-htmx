#!/bin/bash
set -e

echo "Installing dependencies..."
npm ci

echo "Generating templ files..."
templ generate

echo "Building Tailwind CSS..."
npm run build

echo "Building Go application..."
go build -o main .

echo "Build complete!"
