# ao artifact

Read a file from another AO session's artifact directory as raw text. This is
the direct handoff path for architecture RFCs: the orchestrator sends the
architect session id and relative artifact path, and the worker fetches the
full document locally.

## Syntax

```text
ao artifact read --session <architect-session-id> --path <artifact-path>
```

`--path` is relative to the source session's artifact directory. It cannot
escape that directory. The command prints the artifact bytes to stdout and
does not send the content through the orchestrator.

Use this only for artifacts that the orchestrator has explicitly handed off.
The source session id and path are stable AO references; do not guess them.
