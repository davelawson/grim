#!/bin/bash
echo This request should return a user
curl -X GET -H 'Content-Type: application/json' \
  -d '{
    "email": "bob@aol.com"
  }' \
  localhost:8080/user

echo
echo
echo This request should return a 404
curl -X GET -H 'Content-Type: application/json' \
  -d '{
    "email": "bob2@aol.com"
  }' \
  localhost:8080/user
