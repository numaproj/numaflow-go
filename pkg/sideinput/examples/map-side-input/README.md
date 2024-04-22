# Map SideInput Example
An example that demonstrates how to write a [sideinput](https://numaflow.numaproj.io/user-guide/reference/side-inputs/) function along with a sample User Defined function which watches and used the corresponding side input with Mapper function.

### SideInput
```golang
// handle is the side input handler function.
func handle(_ context.Context) sideinputsdk.Message {
	// generate message based on even and odd counter
	counter++
	if counter%2 == 0 {
		return sideinputsdk.BroadcastMessage([]byte("even"))
	}
	// BroadcastMessage() is used to broadcast the message with the given value to other side input vertices.
	// val must be converted to []byte.
	return sideinputsdk.BroadcastMessage([]byte("odd"))
}
```

By using [UDF](#user-defined-function), User perform the retrieval/update for the side input value and perform some operation based via mapper function on it and broadcast/drop a message to other vertices.

### User Defined Function
The UDF vertex will watch for changes to this file and whenever there is a change it will read the file to obtain the new side input value.

### Pipeline spec
In the spec we need to define the side input vertex and the UDF vertex. The UDF vertex will have the side input vertex as a side input.

```yaml
spec:
  sideInputs:
    - name: myticker
      container:
        # A map side input example , see https://github.com/numaproj-contrib/e2e-tests-go/tree/main/map-side-input
        image: "quay.io/numaio/numaproj-contrib/e2e-map-sideinput-example:v0.0.2"
        imagePullPolicy: Always
      trigger:
        schedule: "@every 5s"
```
Vertex spec for the UDF vertex:
```yaml
    - name: si-e2e
      udf:
        container:
          # A map side input udf , see https://github.com/numaproj-contrib/e2e-tests-go/tree/main/map-side-input/udf
          image: "quay.io/numaio/numaproj-contrib/e2e-map-sideinput-udf:v0.0.2"
          imagePullPolicy: Always
      sideInputs:
        - myticker
```



