#!/bin/bash
echo This request should create a user
curl -X POST \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer {$GRIM_TOKEN}'
-d '{
    "email": "bob@aol.com",
    "name": "bob loblaw",
    "password": "password123"
  }' \
  localhost:8080/user
