# dirp

`dirp` は **Directory Processor** の略で、DSL からディレクトリ構造だけを生成する Go 製 CLI です。

## License

This project is licensed under the MIT License.

## 目次 / Table of Contents / Sommaire

- [日本語](#日本語-ja)
- [English](#english-en)
- [Français](#français-fr)

## 日本語 (JA)

### 仕様書

- [日本語仕様書](./dirp_specification_ja.md)
- 他言語: [English](#english-en) / [Français](#français-fr)

### セットアップ（Windows）

1. [Go公式DLページ](https://go.dev/dl/) から `windows-amd64.msi` をインストール
2. 新しい PowerShell で確認

```powershell
go version
gofmt -h
```

3. このプロジェクトでビルド

```powershell
go mod tidy
gofmt -w .\main.go
go build .
```

`.\dirp.exe` が生成されます。

### できること

- 兄弟要素の区切り: `,` / `|` / 改行
- 子階層: `{ ... }`
- テンプレート展開:
  - `#(start,end,step)` 例: `node_#(1,5,2)` -> `node_1,node_3,node_5`
  - `@(a,b,c)` 例: `srv_@(web,db)` -> `srv_web,srv_db`
- Fail-fast な構文エラー通知
- `-f` 指定時は `path:line:col` 形式でエラー表示

### 必要環境

- Go 1.22 以上

### ビルド

```bash
go build .
```

Windows では `dirp.exe` が生成されます。

### 使い方

### 1) インライン文字列

```bash
./dirp -root ./out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

Windows:

```powershell
.\dirp.exe -root .\out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

### 2) `.dirp` ファイル入力

```bash
./dirp -root ./out -f ./sample.dirp
```

Windows:

```powershell
.\dirp.exe -root .\out -f .\sample.dirp
```

### エラー表示（ジャンプしやすい形式）

`-f` で読み込んだ場合、次の形式で表示されます:

```text
sample.dirp:11:7: parse error: expected separator ',' '|' or newline
...該当行...
      ^
```

この `path:line:col` は、多くのエディタ/ターミナルでクリックジャンプ対象として認識されます。

### Git 管理ポリシー

このリポジトリでは以下を通常コミット対象外にしています。

- `*.exe`（ビルド成果物）
- `out/`

### ただし `.exe` を GitHub に載せたい場合

一時的に次で追加できます:

```bash
git add -f dirp.exe
```

または `.gitignore` の `*.exe` を外して運用してください。

## English (EN)

`dirp` stands for **Directory Processor**, a Go CLI that generates directory structures from a DSL.

### Specifications

- [English Specification](./dirp_specification_en.md)
- Other languages: [日本語](#日本語-ja) / [Français](#français-fr)

### Setup (Windows)

1. Install `windows-amd64.msi` from [Go Downloads](https://go.dev/dl/)
2. Verify in a new PowerShell

```powershell
go version
gofmt -h
```

3. Build in this project

```powershell
go mod tidy
gofmt -w .\main.go
go build .
```

`.\dirp.exe` is generated.

### Features

- Sibling separators: `,` / `|` / newline
- Nested hierarchy with `{ ... }`
- Template expansion:
  - `#(start,end,step)` e.g. `node_#(1,5,2)` -> `node_1,node_3,node_5`
  - `@(a,b,c)` e.g. `srv_@(web,db)` -> `srv_web,srv_db`
- Fail-fast syntax errors
- `path:line:col` error output when using `-f`

### Requirements

- Go 1.22+

### Build

```bash
go build .
```

On Windows, this creates `dirp.exe`.

### Usage

#### 1) Inline DSL string

```bash
./dirp -root ./out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

Windows:

```powershell
.\dirp.exe -root .\out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

#### 2) `.dirp` file input

```bash
./dirp -root ./out -f ./sample.dirp
```

Windows:

```powershell
.\dirp.exe -root .\out -f .\sample.dirp
```

### Error Output (Editor-Friendly)

When loaded with `-f`, parse errors are printed as:

```text
sample.dirp:11:7: parse error: expected separator ',' '|' or newline
...offending line...
      ^
```

Many terminals/editors detect `path:line:col` and make it clickable.

### Git Policy

This repository normally ignores:

- `*.exe` (build artifacts)
- `out/`

To include `dirp.exe` temporarily:

```bash
git add -f dirp.exe
```

Or remove `*.exe` from `.gitignore` for your workflow.

## Français (FR)

`dirp` signifie **Directory Processor**, un CLI Go qui génère des structures de répertoires à partir d'un DSL.

### Spécifications

- [Spécification française](./dirp_specification_fr.md)
- Autres langues : [日本語](#日本語-ja) / [English](#english-en)

### Installation (Windows)

1. Installer `windows-amd64.msi` depuis [Go Downloads](https://go.dev/dl/)
2. Vérifier dans un nouveau PowerShell

```powershell
go version
gofmt -h
```

3. Builder dans ce projet

```powershell
go mod tidy
gofmt -w .\main.go
go build .
```

`.\dirp.exe` est généré.

### Fonctionnalités

- Séparateurs de fratrie: `,` / `|` / saut de ligne
- Hiérarchie imbriquée avec `{ ... }`
- Expansion de templates:
  - `#(start,end,step)` ex. `node_#(1,5,2)` -> `node_1,node_3,node_5`
  - `@(a,b,c)` ex. `srv_@(web,db)` -> `srv_web,srv_db`
- Erreurs syntaxiques en mode fail-fast
- Format d'erreur `path:line:col` avec `-f`

### Prérequis

- Go 1.22+

### Build

```bash
go build .
```

Sous Windows, cela produit `dirp.exe`.

### Utilisation

#### 1) Chaîne DSL inline

```bash
./dirp -root ./out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

Windows:

```powershell
.\dirp.exe -root .\out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

#### 2) Entrée fichier `.dirp`

```bash
./dirp -root ./out -f ./sample.dirp
```

Windows:

```powershell
.\dirp.exe -root .\out -f .\sample.dirp
```

### Format d'Erreur (Compatible éditeur)

Avec `-f`, les erreurs sont affichées comme:

```text
sample.dirp:11:7: parse error: expected separator ',' '|' or newline
...ligne concernée...
      ^
```

Le format `path:line:col` est souvent cliquable dans les terminaux/éditeurs.

### Politique Git

Le dépôt ignore normalement:

- `*.exe` (artefacts de build)
- `out/`

Pour inclure `dirp.exe` temporairement:

```bash
git add -f dirp.exe
```

Ou retirez `*.exe` du `.gitignore` selon votre workflow.
