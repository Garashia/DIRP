# Guide d'utilisation detaille de dirp (Français)

Ce document regroupe les usages pratiques trop detailles pour `README.md`.

## 1. Idee principale

- `dirp` analyse un DSL et cree des arborescences de repertoires (dossiers uniquement).
- Separateurs de freres: `,` / `|` / saut de ligne.
- Hierarchie imbriquee avec `{ ... }`.

Exemple:

```text
app{api,web}|docs
```

## 2. Modes d'entree

### 2.1 Entree inline

```bash
./dirp -root ./out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

Windows:

```powershell
.\dirp.exe -root .\out "app{api_#(1,3,1),srv_@(web,db)}|docs"
```

### 2.2 Entree fichier (`-f`)

```bash
./dirp -root ./out -f ./sample.dirp
```

Windows:

```powershell
.\dirp.exe -root .\out -f .\sample.dirp
```

## 3. Expansion de templates

### 3.1 Intervalle numerique `#(start,end,step)`

```text
api_#(1,5,2)
```

Devient:

```text
api_1, api_3, api_5
```

Notes:

- `step` ne peut pas etre 0.
- 3 arguments sont obligatoires.

### 3.2 Liste `@(a,b,c)`

```text
srv_@(web,db,cache)
```

Devient:

```text
srv_web, srv_db, srv_cache
```

## 4. Modes test / debug

### 4.1 Mode sans creation (`-test`)

Affiche l'arbre parse sans creer de dossiers.

```bash
./dirp -c "app{src,bin}" -test
```

### 4.2 Validation en lot (`-cases`)

Un pattern DSL par ligne.

```text
# commentaires ignores
app{src,bin}
service_@(api,web){v_#(1,2,1)}
```

```bash
./dirp -cases ./cases.dirp
```

## 5. Sortie JSON (integration)

Avec `--json`, `dirp` renvoie l'AST ou l'erreur en JSON.

```bash
./dirp -c "app{src,bin}" -test --json
```

En succes, `nodes` utilise les cles `name` / `children`.

## 6. Format d'erreur

Avec `-f`, le format est `path:line:col`.

```text
sample.dirp:11:7: parse error: expected separator ',' '|' or newline
```

## 7. Utilisation comme bibliotheque Go

Vous pouvez utiliser `pkg/dirp` directement:

```go
nodes, err := dirp.Parse(`app{api_#(1,3,1),web}`)
if err != nil {
    err = dirp.WithLineCol(err, `app{api_#(1,3,1),web}`)
}
err = dirp.Build("./out", nodes)
```

## 8. Pieges frequents

- Ne pas combiner `-c` et `-f`.
- Si `-test` et `-mkdir` sont tous les deux presents, `-test` est prioritaire.
- Les caracteres reserves et templates invalides provoquent des erreurs immediates.
