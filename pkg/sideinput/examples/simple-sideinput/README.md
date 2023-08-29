# SideInput Example

An example that demonstrates how to write a `sideinput` function along with a sample `User Defined function` 
which watches and used the corresponding side input.

### SideInput
```golang
// handle is the side input handler function.
func handle(_ context.Context) sideinputsdk.Message {
	t := time.Now()
	// val is the side input message value. This would be the value that the side input vertex receives.
	val := "an example: " + string(t.String())
	// randomly drop side input message. Note that the side input message is not retried.
	// NoBroadcastMessage() is used to drop the message and not to
	// broadcast it to other side input vertices.
	if rand.Int()%2 == 0 {
		return sideinputsdk.NoBroadcastMessage()
	}
	// BroadcastMessage() is used to broadcast the message with the given value to other side input vertices.
	// val must be converted to []byte.
	return sideinputsdk.BroadcastMessage([]byte(val))
}
```
After performing the retrieval/update for the side input value the user can choose to either broadcast the 
message to other side input vertices or drop the message. The side input message is not retried.

For each side input there will be a file with the given path and after any update to the side input value the file will 
be updated.

The directory is fixed and can be accessed through sideinput constants `sideinput.DirPath`.
The file name is the name of the side input.
```golang
sideinput.DirPath -> "/var/numaflow/side-inputs"
sideInputFileName -> "/var/numaflow/side-inputs/sideInputName"
```

### User Defined Function

The UDF vertex will watch for changes to this file and whenever there is a change it will read the file to obtain the new side input value.


### Pipeline spec

In the spec we need to define the side input vertex and the UDF vertex. The UDF vertex will have the side input vertex as a side input.

Side input spec:
```yaml
spec:
  sideInputs:
    - name: myticker
      container:
        image: "quay.io/numaio/numaflow-go/sideinput:v0.5.0"
        imagePullPolicy: Always
      trigger:
        schedule: "*/2 * * * *"

```

Vertex spec for the UDF vertex:
```yaml
    - name: si-log
      udf:
        container:
          image: "quay.io/numaio/numaflow-go/udf-sideinput:v0.5.0"
          imagePullPolicy: Always
      containerTemplate:
        env:
          - name: NUMAFLOW_DEBUG
            value: "true" # DO NOT forget the double quotes!!!
      sideInputs:
        - myticker
```



