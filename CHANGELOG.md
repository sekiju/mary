# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] — Unreleased

### Breaking

- Repository moved from `github.com/sekiju/mdl` to `github.com/knst0/mdl`
- Module path changed to `github.com/knst0/mdl`
- Replaced `sekiju/htt` with `resty.dev/v3` for HTTP client
- Replaced `goccy/go-json` with `bytedance/sonic` for JSON handling
- Refactored config into `FileConfig` + `RuntimeFlags` split
- Threaded `context.Context` through downloader and extractors for cancellation support
- Chapter-level parallelism is now properly configurable via `max_parallel_chapters`

### Added

- TUI mode (`mdl` with no arguments launches interactive terminal UI)
- `manga.ChapterListingFeature` interface for extractor capability detection
- Proper sentinel errors for paid chapters, unsupported chapter listing, etc.
- `check_updates` config option to check GitHub for new releases

### Changed

- `go.mod` targets Go 1.25
- CLI uses `flag.FlagSet` subcommands (`chapters` mode)
- Logger replaced with `rs/zerolog`
- Image decoding supports AVIF and WebP in addition to JPEG/PNG
- Download directory and file format are now configurable via `output {}` block

### Fixed

- Nil-pointer dereference when no cookie configured for cmoa extractor
- `FindChapters` panics replaced with graceful unsupported errors
- Paid chapter detection for ganma extractor
- Per-page download errors now properly surfaced instead of silently swallowed
- Update check no longer blocks download on network failure
