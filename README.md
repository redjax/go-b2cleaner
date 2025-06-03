# B2Cleaner <!-- omit in toc -->

A Backblaze B2 storage cleaning utility written in Go.

## Table of Contents <!-- omit in toc -->

- [Setup](#setup)
  - [Configuration](#configuration)
    - [CLI args](#cli-args)
    - [Environment variables](#environment-variables)
    - [Configuration file](#configuration-file)
- [Usage](#usage)
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

## Usage

Run `b2cleaner --help` to see all options. To see options for a subcommand, run `b2cleaner <command> --help`.

Basic operations:

| Command                                                                                                                                                      | Description                                                                                                                                                       |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `list --bucket <bucket> --path <path/in/bucket> [-c path/to/config.toml] [--sort [name, size]] [--order [asc,desc]]`                                         | List files in a bucket                                                                                                                                            |
| `clean --bucket <bucket> --path <path/in/bucket> [-c path/to/config.toml] [--filetype="fileExt"]+ [-o/--output path/to/results.csv] [--recurse] [--dry-run]` | Clean/delete files in a given bucket/path. Add `--dry-run` to show what will be deleted before deleting. Add `-o/--output` to save deleted objects to a CSV file. |

## Links
