# Batch Map Flatmap

An example User Defined Function that demonstrates how to write a batch map based `flatmap` User Defined Function.


Some considerations for batch map are as follows

- The user will have to ensure that the BatchResponse is tagged with the correct request ID as this will be used by Numaflow for populating information required for system correctness like MessageID for the ISB deduplication.


- The user will have to ensure that all the length of the BatchResponses is equal to the number of requests received. This means that for **each request** there is a BatchResponse.