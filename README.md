# knst0/mdl

> [!IMPORTANT]
> The project is made for educational purposes. If you believe your rights are being violated, contact main contributor.

`mdl` is a terminal manga downloader. Give it a chapter URL, it queues the download in an interactive TUI and saves the
result in the folder.

## Getting started

If you don't want to build the app yourself checkout the [releases page](https://github.com/knst0/mdl/releases).

### Usage

```shell
mdl [OPTIONS] [chapterURL...]
```

Running `mdl` with no arguments opens the TUI, where you can paste URLs, track progress, and manage settings
interactively. Passing one or more chapter URLs on the command line opens the TUI with those chapters already queued.

### Available options

- `--config (string)`: Path to the configuration file. If empty, `MDL_CONFIG` env var is checked, then the OS default
  (`~/.config/mdl/config.json` on Linux, `%USERPROFILE%/mdl/config.json` on Windows,
  `~/Library/Application Support/mdl/config.json` on macOS).
