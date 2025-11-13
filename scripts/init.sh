#!/bin/sh

set -e

mkdir -p /app/run/repo
mkdir -p /app/run/www

KEY_FILE=${GIT_SSH_KEY:-/run/secrets/github_ssh_key}
REPO_URL=${REPO_URL:-git@github.com:kites262/kites262.github.io.git}

GIT_SSH_COMMAND="ssh -i $KEY_FILE -o StrictHostKeyChecking=no" \
  git clone $REPO_URL /app/run/repo
