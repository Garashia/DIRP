# dirp Detailed Usage Guide (English)

This document provides practical usage details that are too verbose for `README.md`.

## 1. Core Concept

- `dirp` parses a DSL and creates directory trees (folders only).
- Sibling separators: `,` / `|` / newline.
- Nested hierarchy uses `{ ... }`.

Example:

```text
app{api,web}|docs
```

## 2. Input Modes

### 2.1 Inline input

```bash
./dirp -root ./out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

Windows:

```powershell
.\dirp.exe -root .\out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

### 2.2 File input (`-f`)

```bash
./dirp -root ./out -f ./sample.dirp
```

Windows:

```powershell
.\dirp.exe -root .\out -f .\sample.dirp
```

## 3. Template Expansion

### 3.1 Numeric range `#(start,end,step)`

```text
api_#(1,5,2)
```

Expands to:

```text
api_1, api_3, api_5
```

Notes:

- `step` must not be 0.
- Exactly 3 arguments are required.

### 3.2 List expansion `@(a,b,c)`

```text
srv_@(web,db,cache)
```

Expands to:

```text
srv_web, srv_db, srv_cache
```

## 4. Test / Debug Modes

### 4.1 Parse-only mode (`-test`)

Shows parsed tree without creating directories.

```bash
./dirp -c "app{src,bin}" -test
```

### 4.2 Batch validation (`-cases`)

Use one DSL pattern per line.

```text
# comments are ignored
app{src,bin}
service_@(api,web){v_#(1,2,1)}
```

```bash
./dirp -cases ./cases.dirp
```

## 5. JSON Output (Integration)

Use `--json` to get AST/error as JSON.

```bash
./dirp -c "app{src,bin}" -test --json
```

On success, `nodes` use `name` / `children` keys.

## 6. Error Format

When using `-f`, errors are reported as `path:line:col`.

```text
sample.dirp:11:7: parse error: expected separator ',' '|' or newline
```

## 7. Using as a Go Library

You can consume `pkg/dirp` directly:

```go
nodes, err := dirp.Parse(`app{api_#(1,3,1),web}`)
if err != nil {
    err = dirp.WithLineCol(err, `app{api_#(1,3,1),web}`)
}
err = dirp.Build("./out", nodes)
```

## 8. Common Pitfalls

- Do not pass `-c` and `-f` together.
- If both `-test` and `-mkdir` are set, `-test` wins.
- Reserved characters or malformed templates cause immediate parse errors.
