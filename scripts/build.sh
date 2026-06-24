#!/usr/bin/env bash

set -e

ROOT_DIR="${BLOG_ROOT:-$(cd "$(dirname "$0")/.." && pwd)}"
APP_DIR="$ROOT_DIR/apps/blog-web"
API_DIR="$ROOT_DIR/apps/admin-api"
RELEASES_DIR="$ROOT_DIR/releases"
CURRENT_LINK="$ROOT_DIR/current"
TIME=$(date +%Y%m%d_%H%M%S)
RELEASE_DIR="$RELEASES_DIR/$TIME"

echo "=== Blog Build Start ==="
echo "Root: $ROOT_DIR"
echo "Release: $RELEASE_DIR"

# Step 1: Export settings from DB to JSON for static site generation
echo "--- Exporting site settings ---"
SITE_CONFIG="$APP_DIR/public/data/site-config.json"
mkdir -p "$(dirname "$SITE_CONFIG")"
"$API_DIR/admin-api" -export-settings "$SITE_CONFIG" 2>/dev/null || echo "Warning: Could not export settings (API DB may not exist yet)"
if [ -f "$SITE_CONFIG" ]; then
    echo "Settings exported: $(cat "$SITE_CONFIG")"
fi

cd "$APP_DIR"

echo "--- Installing dependencies ---"
pnpm install --frozen-lockfile

echo "--- Building Astro ---"
pnpm build

echo "--- Creating release ---"
mkdir -p "$RELEASE_DIR"
cp -r dist/* "$RELEASE_DIR"

# Update current symlink
if [ -L "$CURRENT_LINK" ] || [ -f "$CURRENT_LINK" ]; then
    rm -f "$CURRENT_LINK"
fi
ln -sfn "$RELEASE_DIR" "$CURRENT_LINK"

echo "=== Build Success: $RELEASE_DIR ==="
