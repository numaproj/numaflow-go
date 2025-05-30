# Test Retry Strategy Sink Example

For each message received, sink will create a failure response, containing id in the message.
This sink is used by Numaflow E2E testing for testing exponential backoff retry strategy.
