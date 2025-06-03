# B2Cleaner <!-- omit in toc -->

A Backblaze B2 storage cleaning utility written in Go.

## Table of Contents <!-- omit in toc -->

- [Setup](#setup)
  - [Configuration](#configuration)
    - [CLI args](#cli-args)
    - [Environment variables](#environment-variables)
    - [Configuration file](#configuration-file)
- [Links](#links)

## Setup

### Configuration

Configuration options can be passed to the CLI as `--options`, or loaded from one of the methods below. The sections are arranged in the order they are loaded by the program (CLI flags, environment variables, then config file(s)).

#### CLI args

You can pass environment variables as CLI args by prepending them to a command. For example:

```shell
B2_KEY_ID="your-key-id" \
    B2_APP_KEY="your-app-key" \
    B2_BUCKET="your-bucket" \
    B2_PATH="your/path/" \
    B2_RECURSE="true" \
    b2clean list
```

#### Environment variables

* `B2_KEY_ID="your-key-id"`
* `B2_APP_KEY="your-app-key"`
* `B2_BUCKET="your-bucket"`
* `B2_PATH="your/path/"`
* `B2_RECURSE="true"`

#### Configuration file

Copy the [example `config.toml` file](./example.config.toml) to `config.toml` and set your configuration values:

```toml
key_id  = "your-key-id"
app_key = "your-app-key"
bucket  = "your-bucket"
path    = "your/path/"
recurse = true
```

## Links
