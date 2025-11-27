#!/bin/zsh

# Load API_TOKEN from .env
API_TOKEN=$(grep '^API_TOKEN=' .env | cut -d '=' -f2)

HOST=https://www.emmajelk.se
#HOST=http://localhost:8080

# Download the database files using the API token
curl -H "Authorization: Bearer $API_TOKEN" $HOST/api/db --output ./down/database.db
curl -H "Authorization: Bearer $API_TOKEN" $HOST/api/db-shm --output ./down/database.db-shm
curl -H "Authorization: Bearer $API_TOKEN" $HOST/api/db-wal --output ./down/database.db-wal
