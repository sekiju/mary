# knst0/mdl

> [!IMPORTANT]
> The project is made for educational purposes. If you believe your rights are being violated, contact main contributor.

## Getting started

If you don't want to build the app yourself checkout the [releases page](https://github.com/knst0/mdl/releases).

### Usage

```shell
mdl [OPTIONS] chapterURL [chapterURLs...]
```

### Download example

```shell
mdl https://comic-ogyaaa.com/episode/4856001361536369722 https://comic-ogyaaa.com/episode/4856001361561258284
```

### Available options

- `--cookie (string)`: Provides the cookie string for the current session
- `--config (string)`: Path to the configuration file. If empty, `MDL_CONFIG` env var is checked, then the OS default (`~/.config/mdl/config.json` on Linux, `%AppData%/mdl/config.json` on Windows, `~/Library/Application Support/mdl/config.json` on macOS).

## Config

The config file is JSON and lives in an OS-appropriate location. On first run, a default config is automatically generated. The format is documented by the JSON Schema at [`schema/config.schema.json`](schema/config.schema.json).

Key settings:

- `application.max_parallel_downloads` — concurrent page/image downloads (default `4`)
- `application.max_parallel_chapters` — concurrent chapter downloads (default `1`)
- `output.directory` — download output directory (default `downloads`)
- `output.file_format` — output format: `auto`, `png`, `jpeg`, `avif`, or `webp`
- `site.<hostname>.cookie` — per-site session cookie for authenticated access

## Architecture

Site extractors live under `extractor/` and implement `sdk/manga.Extractor`.
They self-register via `extractor/registry.Register` in `init()` functions,
so the core package never imports concrete extractor packages by name.
Built-in extractors are wired in through blank imports in
`extractor/builtin.go`.

See [Writing a site extractor](docs/extractors.md) for a walkthrough.

## Development

### Build

```shell
go build ./...
go vet ./...
go test ./...
```

### Releasing

To cut a release, run `make release-patch` (or `-minor`/`-major`).
This merges changesets, bumps the version, tags, and pushes.
