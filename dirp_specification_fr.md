# Spécification Technique du DSL : dirp (Directory Processor)

**Version :** 1.0  
**Extension de fichier :** `.dirp`  
**Encodage :** UTF-8 (recommandé)

---

## 1. Introduction
`dirp` est un langage spécifique au domaine (DSL) conçu exclusivement pour la création de structures de répertoires. Sa philosophie repose sur le minimalisme, la prévisibilité et une gestion stricte des erreurs. Contrairement à d'autres outils, `dirp` ne crée jamais de fichiers ; il définit uniquement des conteneurs (dossiers).

## 2. Syntaxe de Base

### 2.1. Entités
Toute chaîne de caractères non réservée est interprétée comme un nom de répertoire.

### 2.2. Séparateurs de Fratrie
Les répertoires situés au même niveau hiérarchique sont séparés par l'un des trois caractères suivants, traités de manière identique :
- La virgule `,`
- La barre verticale `|`
- Le saut de ligne `\n`

### 2.3. Hiérarchie (Parent-Enfant)
L'imbrication est définie par des accolades `{ }`. 
- `A { B }` signifie que le répertoire `B` est créé à l'intérieur du répertoire `A`.
- Les accolades peuvent être imbriquées sans limite théorique de profondeur.

### 2.4. Gestion des Espaces
- Les espaces situés au début et à la fin d'un nom de répertoire sont ignorés (trimming).
- Les espaces situés à l'intérieur d'un nom sont préservés et font partie intégrante du nom.

---

## 3. Fonctions de Gabarit (Templating)

Une entité peut contenir au maximum **une seule** fonction de gabarit. L'utilisation récursive ou multiple de fonctions au sein d'un même nom est strictement interdite.

### 3.1. Fonction de Plage : `#(début, fin, pas)`
Génère une séquence numérique.
- **Paramètres :** Doivent être des entiers.
- **Exemple :** `node_#(1, 5, 2)` génère `node_1`, `node_3`, `node_5`.

### 3.2. Fonction de Liste : `@(item1, item2, ...)`
Génère une série de répertoires basés sur les éléments listés.
- **Paramètres :** Chaînes de caractères séparées par des virgules. Les espaces entourant les items sont ignorés.
- **Exemple :** `srv_@(web, db)` génère `srv_web`, `srv_db`.

### 3.3. Contraintes des Fonctions
- **Séparateur interne :** Seule la virgule `,` est autorisée à l'intérieur des parenthèses.
- **Sauts de ligne :** Autorisés entre les arguments (après une virgule), mais strictement interdits à l'intérieur d'une valeur ou d'un nombre.
- **Unicité :** `nom_#(1,3)_@(a,b)` est syntaxiquement invalide.

---

## 4. Gestion des Erreurs (Fail-Fast Policy)

Le processeur `dirp` doit interrompre immédiatement l'exécution et signaler une erreur détaillée dans les cas suivants :

1. **Nom Vide :** Présence de séparateurs consécutifs (ex: `A,,B` ou `{,}`) sans nom défini.
2. **Caractères Réservés :** Utilisation de `{ } ( ) , | # @` à l'intérieur d'un nom de répertoire en dehors de leur fonction syntaxique.
3. **Déséquilibre :** Accolades ou parenthèses non fermées.
4. **Syntaxe de Fonction Invalide :** Mauvais nombre d'arguments ou types de données incorrects dans `#( )` ou `@( )`.
5. **Recursion de Gabarit :** Tentative d'imbriquer une fonction dans une autre.

---

## 5. Comportement du Système

### 5.1. Existence Préalable
Si un répertoire cible existe déjà, `dirp` ignore l'étape de création pour cette entité spécifique et poursuit son exécution sans erreur. Aucune modification n'est apportée aux répertoires existants.

### 5.2. Caractères Spéciaux et Encodage
- Bien que l'UTF-8 soit recommandé, le support des caractères non-ASCII (ex: Japonais, accents français) dépend du système de fichiers de l'OS hôte.
- `dirp` ne garantit pas la compatibilité des noms de fichiers entre différents OS si des caractères spéciaux sont utilisés.
- Les symboles de fonction (ex: `#`, `@`) doivent obligatoirement être en demi-chasse (ASCII).

### 5.3. Exécution par une Seule Ligne
Le langage est conçu pour être condensable en une seule ligne sans perte de structure, facilitant l'intégration via CLI ou scripts tiers.

---

## 6. Philosophie d'Extension
Les logiques complexes (expressions régulières, calculs dynamiques) ne doivent pas être intégrées au cœur de `dirp`. Elles doivent être traitées par un langage de script parent (tel que Python) qui générera ensuite la chaîne de caractères conforme à la présente spécification.
