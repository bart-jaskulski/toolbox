# toolbox

Monorepo for mini Go CLI tools

## Commands

- `bm`: bookmark manager — `go install github.com/bart-jaskulski/toolbox/bm@latest`
- `commitment`: AI-assisted commit message generator — `go install github.com/bart-jaskulski/toolbox/commitment@latest`
- `ddc`: DevDocs offline browser — `go install github.com/bart-jaskulski/toolbox/ddc@latest`
- `em`: terminal emoji picker — `go install github.com/bart-jaskulski/toolbox/em@latest`
- `mfeed`: RSS/Atom meta-feed generator — `go install github.com/bart-jaskulski/toolbox/mfeed@latest`
- `rdr`: readability-based web page to Markdown reader — `go install github.com/bart-jaskulski/toolbox/rdr@latest`

## Build

Compile any tool locally without installing (avoids clashing with the module directories):

```bash
mkdir -p bin
go work sync
go build -o bin/bm ./bm          # change bm -> commitment|ddc|em|mfeed|rdr
```

## Local development

The repo uses a Go workspace (`go.work`) so running `go build` or `go test` from any tool will use local sources:

```bash
go work sync
go test ./...
```
