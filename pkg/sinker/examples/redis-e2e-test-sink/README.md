# Redis E2E Test Sink
A User Defined Sink using redis hashes to store messages.
The name of a hash is pipelineName:sinkName.

For each message received, the sink will store the message in a hash with the key being the payload of the message
and the value being the no. of occurrences of that payload so far.

This sink is used by Numaflow E2E testing.

