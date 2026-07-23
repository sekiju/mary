# sekiju/mdl

> [!IMPORTANT]
> The project is made for educational purposes. If you believe your rights are being violated, contact main contributor.

## Getting started

If you don't want to build the app yourself checkout the [releases page](https://github.com/sekiju/mdl/releases).

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
- `--config (string)`: Specifies the path to the configuration file (default: `config.hcl`)

## Config

Refer to [example.config.hcl](example.config.hcl) for the default configuration settings.

### Optimizing download threads for best performance

On average, each image is approximately 400 KB in size. The `download.concurrent_processes` setting determines the number of simultaneous download threads. It’s important to note that increasing this value beyond `internet speed / image size` offers no additional benefit.

Using only one thread (`1`) is generally inefficient compared to multiple threads. For instance, downloading a chapter with `4` threads might take 25 seconds, whereas a single thread could require up to 42 seconds.

#### Example Scenario
- **Internet Speed**: 100 Mbps
- **Threads**:
    - `4 threads`: Download completed in **1m25s**
    - `12 threads`: Download completed in **38s**
    - `50 threads`: Download completed in **34s**

From this, the optimal number of threads is approximately **13**.

#### Why not always use maximum threads?
While increasing threads can improve download speeds, setting too high a number, like `50`, places unnecessary strain on your computer’s resources, including disk, processor, and memory usage. This extra load does not significantly improve performance and may degrade overall system efficiency. Therefore, it’s advisable to determine and use an optimal thread count for your setup.

### Paid chapter download or how to obtain `cookie`?

To extract cookies from websites, use tools like [Cookie-Editor](https://cookie-editor.com).

Steps:

1. Open the extension on the target website.
2. Click the Export button (located in the bottom-right corner).
3. Select the Header string option.

Now in clipboard you have Cookie Header string, paste it to config or use with `--cookie` CLI argument.

This copies the Cookie Header string to your clipboard. You can then paste it into the configuration file or provide it via the `--cookie`
CLI argument or modify config file:

```hcl
// ...previous config

site {
  "shonenjumpplus.com" {
    cookie = "glsc=1hYa4GrNp2DndSNIShVyoDGP6MgDmaJhiX22C0X734hkzod56wsBN7Fy1S5ZBOQd"
  }
}
```

or with CLI option example:

```shell
mdl --cookie glsc=1hYa4GrNp2DndSNIShVyoDGP6MgDmaJhiX22C0X734hkzod56wsBN7Fy1S5ZBOQd https://comic-ogyaaa.com/episode/4856001361561258284
```

## Development

### Build, vet, test

```shell
go build ./...
go vet ./...
go test ./...                    # unit tests only, no network required
go test -tags=integration ./...  # also exercises live-site extractor paths, requires network
```

### Lint

```shell
golangci-lint run
```

Configuration lives in [.golangci.yaml](.golangci.yaml).

### Architecture overview

- **`extractor/`** — one package per supported site, each implementing the
  `sdk/manga.Extractor` interface (`FindChapters`, `FindChapter`,
  `FindChapterPages`, `SetSettings`). `extractor/extractor.go` holds a
  `domainRegistry map[string]Factory` mapping exact hostnames to extractor
  constructors (`getSession` does exact-match lookup, no subdomain
  fallback). Extractors that share a page format (e.g. Speed Binb) embed a
  common `extractor/template/*` package; all extractors embed a shared
  `base`/`util.Base` for `settings`/`SetSettings` and cookie-header
  injection boilerplate.
- **`downloader/`** — `Downloader.Queue` pushes work onto a channel
  consumed by `run()`, which resolves the extractor, finds chapter pages,
  and fans out page downloads under a `MaxParallelDownloads`-bounded
  semaphore; chapter-level concurrency is bounded separately by
  `MaxParallelChapters` (both configured under `application` in
  `config.hcl`, default sequential when unset). A `context.Context` is
  threaded through the whole path so `Ctrl+C` cancels in-flight downloads
  and cleans up partially-written chapter directories.
- **`config/`** — `config.Params` is a package-level singleton split into
  `FileConfig` (everything unmarshaled from `config.hcl` via `koanf`:
  `Application`, `Output`, `Sites`) and `RuntimeFlags` (CLI-only state such
  as `PrimaryCookie`/`ListChaptersMode`/`DownloadChapters`, populated by
  `cmd/cli/main.go` flag parsing, never touched by `koanf.Unmarshal`).
- **`sdk/manga/`** — the extractor-facing public types/interfaces
  (`Extractor`, `Chapter`, `Page`, `Settings`) and the `Err*` sentinel error
  vocabulary extractors are expected to return for recognizable domain
  conditions (not-found, paid, unsupported, credentials-required, etc.)
  instead of ad hoc errors or panics.