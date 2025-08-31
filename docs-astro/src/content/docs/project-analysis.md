---
title: Project Analysis
description: Language detection, dependency extraction, and entry point identification
---

## Language Detection

Automatically identifies project language and framework:

- **Go** — Detects modules, entry points, test patterns
- **JavaScript/TypeScript** — Node.js, React, Vue, Angular projects
- **Python** — Packages, virtual environments, requirements
- **Rust** — Cargo workspaces, crates, features
- **Java/Kotlin** — Maven, Gradle, Spring projects

## File Classification

Intelligent categorization:

- **Entry Points** — main.go, index.js, app.py, main.rs
- **Configuration** — package.json, go.mod, Cargo.toml, pom.xml
- **Source Code** — Implementation files, modules, packages
- **Tests** — Unit, integration, and e2e test files
- **Documentation** — README, docs, comments, API specs

## Git Integration

Repository context:

```bash
# Example git information included
Branch: main
Commit: a7d7640 - "feat: Enhance .gitignore patterns"
Status: 2 modified, 1 untracked
```

## Dependencies

Extracts dependency information from:

| Language | File | Details |
|----------|------|----------|
| Go | `go.mod` | Modules, versions, replace directives |
| Node.js | `package.json` | Dependencies, scripts, engines |
| Python | `requirements.txt` | Packages, version constraints |
| Rust | `Cargo.toml` | Crates, features, workspaces |
| Java | `pom.xml` | Dependencies, plugins, profiles |
