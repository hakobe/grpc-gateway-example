#!/bin/bash
trap 'kill $(jobs -p)' EXIT

./server/server &
./gateway/gateway &

sleep 3

curl -X POST -H"Content-Type: application/json" \
    --data-binary '{ "article": {"title": "hello1", "body": "Hello, World!", "created":"2009-11-10T23:00:00Z"}}' \
    http://localhost:5050/articles/post

curl -X POST -H"Content-Type: application/json" \
    --data-binary '{ "article": {"title": "hello2", "body": "Hello, World!!", "created":"2018-04-18T02:06:00Z"}}' \
    http://localhost:5050/articles/post

curl -X GET http://localhost:5050/articles/recent
