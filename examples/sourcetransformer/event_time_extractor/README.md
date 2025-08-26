# Event Time Extractor


A eventTimeExtractor transformer extracts event time from the payload of the message, based on a user-provided `expression` and an optional `format` specification.

- `expression` is used to compile the payload to a string representation of the event time.
- `format` is used to convert the event time in string format to a time.Time object.


In this example, we have a `expression`, which is used to compile the time of second item in the json payload to a string.