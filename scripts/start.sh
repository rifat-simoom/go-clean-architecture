#!/bin/bash

# Load environment variables from .env file
cd /var/www/html/repos/go-hexarch/
firebase emulators:start --only firestore &
echo "Waiting for Firestore to start..."
while ! nc -z localhost 8787; do
    sleep 2
done
