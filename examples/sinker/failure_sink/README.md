# Failure Sink Example

For each message received, sink will create a `failure` response, containing id in the message.
The [Retry Strategy](https://numaflow.numaproj.io/user-guide/sinks/retry-strategy/) decides how Numaflow responds to the failed messages.