#EvenOdd

An example User Defined Function that checks if the input message is an integer, and returns it with `even` or `odd` key, otherwise drops it.

This can be used to demonstrate how to do [Conditional Forwarding](https://github.com/numaproj/numaflow/blob/main/docs/CONDITIONAL_FORWARDING.md).

```yaml
apiVersion: numaflow.numaproj.io/v1alpha1
kind: Pipeline
metadata:
  name: even-odd
spec:
  vertices:
    - name: in
      source:
        http: {}
    - name: even-or-odd
      udf:
        container:
          image: quay.io/numaio/go-even-odd-example
    - name: even-sink
      sink:
        log: {}
    - name: odd-sink
      sink:
        log: {}
  edges:
    - from: in
      to: even-or-odd
    - from: even-or-odd
      to: even-sink
      conditions:
        keyIn:
          - even
    - from: even-or-odd
      to: odd-sink
      conditions:
        keyIn:
          - odd
```
