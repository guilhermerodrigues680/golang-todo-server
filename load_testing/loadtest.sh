#!/bin/bash

# bash loadtest.sh

NUM_REQUESTS='10000'
NUM_CONCURRENCY='150'
AUTH_HEADER='Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTgzNTE4NTgsImlhdCI6MTYxODM0ODI1OCwiaXNzIjoidG9kby1zZXJ2ZXIiLCJzdWIiOiJndWkifQ.pO4HyBeRot-asp_8VwB2fMDquTOFauLA6rWKIaWdfu4'

ab -c $NUM_CONCURRENCY \
    -n $NUM_REQUESTS \
    -H "${AUTH_HEADER}" \
    -T 'application/json' \
    -p todo.json \
    -e report.csv \
    -w \
    https://todo-service.guilhermerodri8.repl.co/api/v1/todo > result.html
