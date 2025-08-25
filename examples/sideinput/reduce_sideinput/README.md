# Reduce SideInput Example
An example that demonstrates how to write a [sideinput](https://numaflow.numaproj.io/user-guide/reference/side-inputs/) function along with a sample User Defined function which watches and used the corresponding side input with Mapper function.

### SideInput
```golang
// handle is the side input handler function.
func handle(_ context.Context) sideinputsdk.Message {
    counter++
    // BroadcastMessage() is used to broadcast the message with the given value to other side input vertices.
    // val must be converted to []byte.
    return sideinputsdk.BroadcastMessage([]byte(strconv.Itoa(counter)))
}
```

By using [UDF](#user-defined-function), User perform the retrieval/update for the side input value and perform some operation based via reducer function on it and broadcast/drop a message to other vertices.

### User Defined Function
The UDF vertex will watch for changes to this file and whenever there is a change it will read the file to obtain the new side input value.

### Pipeline spec
In the spec we need to define the side input vertex and the UDF vertex. The UDF vertex will have the side input vertex as a side input.

```yaml
spec:
  sideInputs:
    - name: myticker
      container:
        image: "quay.io/numaio/numaflow-go/reduce-sideinput:stable"
        imagePullPolicy: Always
      trigger:
        schedule: "@every 5s"
```
Vertex spec for the UDF vertex:
```yaml
    - name: si-e2e
      udf:
        container:
          image: "quay.io/numaio/numaflow-go/reduce-sideinput-udf:stable"
          imagePullPolicy: Always
        groupBy:
          window:
            fixed:
              length: 10s
          keyed: true
          storage:
            persistentVolumeClaim:
              volumeSize: 10Gi
              accessMode: ReadWriteOnce
      sideInputs:
        - myticker
```



