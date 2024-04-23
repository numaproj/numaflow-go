# Redis E2E Test Sink

This Redis UDSink was specifically created for Numaflow end-to-end tests. 
It interacts with a Redis instance and, upon receiving data, 
it sets the content of the side input as the key in Redis, 
and the number of occurrences of that key as its value.
The name of a hash is pipelineName:sinkName.

# How It Works:
Uses the fsnotify package to watch for changes to the side input file.
Reads and stores the content of the side input file when created.
For each received message, the sink:
1. Uses the content of the side input as the Redis hash field (key)
2. Increments the value associated with that key (number of occurrences of that key)
3. The hash name in Redis is determined by the environment variables `NUMAFLOW_PIPELINE_NAME` and `NUMAFLOW_VERTEX_NAME`


