# EOT Sink

A user-defined sink that works normally at first, then produces an EOT failure after a configurable delay or message count. Intended for progressive-rollout / critical-error AnalysisTemplate e2e tests.

After the trigger, each batch is drained without responses so numaflow records `eot_received_from_sink`.

## Configuration

| Env var | Default | Description |
| --- | --- | --- |
| `FAIL_AFTER_SECONDS` | `60` | Start failing after this many seconds. Set to `0` to disable. |
| `FAIL_AFTER_MESSAGES` | unset / `0` | Start failing after this many successful messages. When both are set, failure starts when either condition is met. |

## Image

`quay.io/numaio/numaflow-go/sink-eot:stable`
