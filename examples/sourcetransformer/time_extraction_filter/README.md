# Time Extraction Filter 


It evaluates a message on a pipeline and if valid, extracts event time from the payload of the messsage.
- `filterExpr` is used to evaluate and drop invalid messages.
- `eventTimeExpr` is used to compile the payload to a string representation of the event time.
- `format` is used to convert the event time in string format to a time.Time object.

For more details, check [builtin timeExtractionFilter.](https://numaflow.numaproj.io/user-guide/sources/transformer/builtin-transformers/time-extraction-filter/).

In this example we have used `filterExpr` and `eventTimeExpr`. 

- Selects json payloads with id < 100 using `filterExpr`
- Selects time field of json payload using `eventTimeExpr` for a string representation of the event time.