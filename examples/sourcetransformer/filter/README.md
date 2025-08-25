# Filter

An example User Defined Function that filter the messages based on expression implemented with `expr` and `sprig` libraries. For more details, check [builtin filter.](https://numaflow.numaproj.io/user-guide/sources/transformer/builtin-transformers/filter/)


In this example, the filter is designed to select json payloads where:

- The `id` is less than 100.
- The `msg` is 'hello'.
- The `desc` contains 'good'.

> Note: We already have a builtin for transformer `filter` in numaflow Go, but not in Rust.