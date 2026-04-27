# dirp DSL Technical Specification (English)

**Version:** 1.0  
**File extension:** `.dirp`  
**Encoding:** UTF-8 (recommended)

---

## 1. Introduction

`dirp` is a domain-specific language focused only on directory-tree creation.  
Its core principles are minimal syntax, predictability, and strict error handling.  
`dirp` never creates files; it only creates directories.

## 2. Core Syntax

### 2.1 Entities
Any non-reserved string is treated as a directory name.

### 2.2 Sibling separators
Directories on the same level are separated by one of:
- `,`
- `|`
- newline (`\n`)

All are treated equivalently.

### 2.3 Parent-child hierarchy
Nested structure is represented by `{ ... }`.
- `A { B }` means create `B` inside `A`
- Nesting depth is theoretically unbounded

### 2.4 Whitespace handling
- Leading/trailing spaces in a directory name are trimmed
- Internal spaces are preserved

---

## 3. Template Functions

An entity name can contain **at most one** template function.  
Recursive or multiple template usage in a single entity is invalid.

### 3.1 Range function `#(...)`
Generates a numeric sequence. It accepts **1 to 3 arguments**.
- `#(end)` -> `start=1`, `step` inferred toward `end`
- `#(start,end)` -> `step` inferred toward `end`
- `#(start,end,step)` -> explicit `step`
- Example: `node_#(3)` -> `node_1`, `node_2`, `node_3`
- Example: `node_#(3,1)` -> `node_3`, `node_2`, `node_1`
- Example: `node_#(1,5,2)` -> `node_1`, `node_3`, `node_5`

### 3.2 List function `@(item1,item2,...)`
Generates names based on listed items.
- Example: `srv_@(web,db)` -> `srv_web`, `srv_db`

### 3.3 Constraints
- Only `,` is allowed as internal separator in function args
- Newlines between arguments are allowed; newline inside a value/number is not
- Leading list template is valid (e.g. `@(a,b)`)
- Leading range template is invalid (`#(...)` must follow a prefix string)
- `name_#(1,3)_@(a,b)` is syntactically invalid

---

## 4. Error Handling (Fail-Fast)

Processing must stop immediately for:
1. Empty name (e.g. `A,,B` or `{,}`)
2. Illegal reserved-character usage (`{ } ( ) , | # @`)
3. Unbalanced braces or parentheses
4. Invalid function argument count/type (`#(...)` accepts 1 to 3 args)
5. Recursive template nesting attempts

---

## 5. Runtime Behavior

### 5.1 Existing directories
If a target directory already exists, creation is skipped and execution continues without error.

### 5.2 Special characters and portability
- UTF-8 is recommended
- Non-ASCII support depends on host OS/filesystem
- Function symbols (`#`, `@`) must be ASCII half-width

### 5.3 One-line compatibility
The language is designed to be representable in one line without semantic loss.

---

## 6. Extension philosophy

Complex logic (regex, dynamic calculations, etc.) should stay outside core `dirp`.  
Such logic should be handled by a parent scripting language that generates valid `dirp` strings.
