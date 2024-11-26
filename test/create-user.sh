#!/bin/bash
echo This request should create a user
curl -X POST -H 'Content-Type: application/json' \
  -d '{
    "email": "new_bob@aol.com",
    "name": "bob loblaw",
    "password_hash": "asdf"
  }' \
  localhost:8080/user
