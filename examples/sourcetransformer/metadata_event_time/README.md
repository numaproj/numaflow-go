# Metadata Event Time

This example demonstrates a Source Transformer that updates the event time of each message to the current time. It also adds this new timestamp to the message's metadata.

## Metadata Details

Each message includes the following user metadata:
- Group: `event-time-group`
- Key: `event-time`
- Value: The new event time in RFC3339 format (e.g., `2023-10-27T10:00:00Z`)
