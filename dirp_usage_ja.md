# dirp 詳細ガイド（日本語）

このドキュメントは、`README.md` に収まりきらない実運用向けの使い方をまとめた詳細版です。

## 1. 基本の考え方

- `dirp` は DSL 文字列を解析し、ディレクトリ構造（フォルダのみ）を作成します。
- 区切りは `,` / `|` / 改行 を使えます。
- 子階層は `{ ... }` で表現します。

例:

```text
app{api,web}|docs
```

## 2. 入力方法

### 2.1 インライン入力

```bash
./dirp -root ./out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

Windows:

```powershell
.\dirp.exe -root .\out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

### 2.2 ファイル入力（`-f`）

```bash
./dirp -root ./out -f ./sample.dirp
```

Windows:

```powershell
.\dirp.exe -root .\out -f .\sample.dirp
```

## 3. テンプレート展開

### 3.1 数値レンジ `#(start,end,step)`

```text
api_#(1,5,2)
```

展開:

```text
api_1, api_3, api_5
```

注意:

- `step` は 0 を許可しません。
- 引数は 3 つ必須です。

### 3.2 リスト展開 `@(a,b,c)`

```text
srv_@(web,db,cache)
```

展開:

```text
srv_web, srv_db, srv_cache
```

## 4. テスト/デバッグ向けモード

### 4.1 非生成テスト（`-test`）

パース結果を表示し、実際のフォルダは作りません。

```bash
./dirp -c "app{src,bin}" -test
```

### 4.2 ケース一括検証（`-cases`）

1行1パターンのファイルを使って連続検証します。

```text
# comments are ignored
app{src,bin}
service_@(api,web){v_#(1,2,1)}
```

```bash
./dirp -cases ./cases.dirp
```

## 5. JSON出力（外部連携向け）

`--json` を指定すると、ASTまたはエラーを JSON で返します。

```bash
./dirp -c "app{src,bin}" -test --json
```

成功時の `nodes` は `name` / `children` キーです。

## 6. エラー表示

`-f` でファイル入力した場合は `path:line:col` 形式で表示されます。

```text
sample.dirp:11:7: parse error: expected separator ',' '|' or newline
```

## 7. APIとして使う（Go）

`pkg/dirp` は CLI 以外からも利用できます。

```go
nodes, err := dirp.Parse(`app{api_#(1,3,1),web}`)
if err != nil {
    // 必要なら line/col を付与
    err = dirp.WithLineCol(err, `app{api_#(1,3,1),web}`)
}
err = dirp.Build("./out", nodes)
```

## 8. よくあるハマりどころ

- `-c` と `-f` の同時指定はできません。
- `-test` と `-mkdir` を同時指定した場合、`-test` が優先されます。
- 全角記号や予約文字の混入で構文エラーになることがあります。
