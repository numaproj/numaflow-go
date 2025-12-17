# Simple Source With Metadata

This example demonstrates a User Defined Source that generates a limited stream of messages (100 total) and attaches metadata to each message.

## Metadata Details

Each message includes the following user metadata:
- Group: `simple-source`
- Key: `txn-id`
- Value: A random UUID (e.g., `550e8400-e29b-41d4-a716-446655440000`)
