# Crash-at-frequency sink example

This sink returns **success** for every message until **wall-clock time** since the last panic (or since the first processed message) reaches `CRASH_INTERVAL`, then it **`panic()`s**. The Go SDK turns that into an internal handler error (`UDF_EXECUTION_ERROR`), which is useful for chaos or resilience testing.

Set the interval with [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) syntax:

- `CRASH_INTERVAL` — default `5m` if unset. Must be positive; invalid values cause the process to exit at startup.

See also the [retry strategy](https://numaflow.numaproj.io/user-guide/sinks/retry-strategy/) docs for how Numaflow handles failures.
