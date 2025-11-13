#!/bin/sh

set -e

KEY_FILE=${GIT_SSH_KEY:-/run/secrets/github_ssh_key}

echo "====== Build Started: $(date) ======"

# Ensure key file exists
if ! cd /app/run/repo; then
  echo "!!!!! Error: Could not change directory to /app/run/repo !!!!!"
  exit 1
fi

# Ensure repo directory exists
cd /app/run/repo || {
  echo "!!!!! Error: Could not change directory to /app/run/repo !!!!!"
  exit 1
}

echo "Pulling latest changes from $WEBHOOK_BRANCH branch..."

# use SSH with your secret key
GIT_SSH_COMMAND="ssh -i $KEY_FILE -o StrictHostKeyChecking=no" \
  git pull origin $WEBHOOK_BRANCH

echo "Running Hugo to build site..."
/usr/local/bin/hugo --destination /app/run/www --minify

echo "✓✓✓✓✓ Build completed: $(date) ✓✓✓✓✓"
echo "========================================"
