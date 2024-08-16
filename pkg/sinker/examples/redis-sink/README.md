# Redis E2E Test Sink
A User Defined Sink using redis hashes to store messages.
The hash key is set by an environment variable `SINK_HASH_KEY` under the sink container spec.

For each message received, the sink will store the message in the hash with the key being the payload of the message
and the value being the no. of occurrences of that payload so far.

This sink is used by Numaflow E2E testing.

