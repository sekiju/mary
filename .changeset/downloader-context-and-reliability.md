---
"mdl": major
---

Threaded `context.Context` through the downloader and extractors for cancellation support. Chapter-level parallelism is now properly configurable via `application.max_parallel_chapters`, and download directory / output file format are configurable via the `output {}` config block. Image decoding now supports AVIF and WebP in addition to JPEG/PNG.

Fixed a panic on double-close of the downloader (`sync.Once` guard in `Downloader.Stop()`), and per-page download errors are now properly surfaced instead of silently swallowed. Logger replaced with `rs/zerolog`. CLI now uses `flag.FlagSet` subcommands (`chapters` mode).
