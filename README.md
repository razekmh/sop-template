# SOP Template

A white-label Standard Operating Procedure template. Fill in your organisation's details, render, and publish a static site via GitHub Actions + Hugo.

---

## How it works

```
src/          ← Markdown source files with {{placeholders}}
     ↓  render (Go binary)
hugo/content/ ← Rendered Markdown with real values
     ↓  Hugo build
hugo/public/  ← Static site, deployed to GitHub Pages
```

`src/` is **never** deployed. Only `hugo/content/` feeds the site.

---

## Quickstart (instance owners)

### 1. Create your repo from this template

Click **Use this template** on GitHub. Do **not** fork — use the template button.

### 2. Clone and set up upstream remote

```bash
git clone https://github.com/YOUR-ORG/YOUR-SOP
cd YOUR-SOP
git remote add upstream https://github.com/SOURCE-ORG/sop-template
```

### 3. Configure your instance

```bash
cp config.example.yml config.yml
```

Edit `config.yml` with your organisation's real values.

### 4. Render

```bash
# macOS (Intel)
./bin/render-mac

# macOS (Apple Silicon)
./bin/render-mac-arm

# Linux
./bin/render-linux

# Windows (Git Bash)
./bin/render-windows.exe
```

This populates `hugo/content/` with your rendered documents.

### 5. Push

```bash
git add hugo/content/ hugo/config.yml
git commit -m "init: render SOP for YOUR ORG NAME"
git push
```

The GitHub Action will build and deploy your site automatically.

---

## Pulling upstream updates

When the upstream template publishes updates:

```bash
git fetch upstream
git merge upstream/main
# Resolve any content conflicts
./bin/render-mac   # re-render after merge
git push
```

Merge conflicts will only appear on genuine content changes — never on substituted values, because the source always uses `{{placeholders}}`.

---

## Repo structure

```
sop-template/
├── .github/
│   └── workflows/
│       ├── lint.yml             # blocks hardcoded values in src/
│       ├── render-and-build.yml # renders + builds + deploys Hugo site
│       └── release.yml          # builds Go binaries on version tag
├── bin/                         # pre-compiled render binaries
├── src/                         # source .md files with {{placeholders}}
├── hugo/
│   ├── config.yml               # Hugo config (also uses placeholders)
│   ├── content/                 # rendered output (populated by render step)
│   ├── static/                  # static assets (logo etc.)
│   └── themes/geekdoc/          # Hugo theme (git submodule)
├── tools/
│   └── render.go                # renderer source
├── config.example.yml           # placeholder reference
├── config.yml                   # YOUR values — gitignored, never commit
└── .gitignore
```

---

## Adding new placeholders

1. Add the placeholder `{{your_key}}` in any `src/*.md` file
2. Add `your_key` to `config.example.yml` with an example value
3. Notify instance owners to add the key to their `config.yml` before re-rendering

The render binary will exit with an error if any placeholder is unresolved.

---

## Building the render binary yourself

Requires Go 1.21+:

```bash
cd tools
go mod tidy
GOOS=linux   GOARCH=amd64 go build -o ../bin/render-linux    render.go
GOOS=darwin  GOARCH=amd64 go build -o ../bin/render-mac      render.go
GOOS=darwin  GOARCH=arm64 go build -o ../bin/render-mac-arm  render.go
GOOS=windows GOARCH=amd64 go build -o ../bin/render-windows.exe render.go
```
