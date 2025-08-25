# Retry

This is a User Defined Function example which imitates a need to retry a certain message some of the time, due to failure
or some other reason. 
For each unique message received, it retries the message twice and then succeeds on the third time.
It does a retry by publishing the message with a "retry" tag (which by the definition of the Pipeline should cause the message to be cycled back to it).

Note that this somewhat silly example would only be used with a single Vertex replica since it makes use of an in-memory map.
