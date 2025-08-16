<h1 align="center">envprof</h1>

<p align="center">
  <img alt="envprof logo" src="assets/images/envprof.png" height="150" />
  <p align="center">Profile-based environment variable manager</p>
</p>

---

[![GitHub release](https://img.shields.io/github/v/release/idelchi/envprof)](https://github.com/idelchi/envprof/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/envprof.svg)](https://pkg.go.dev/github.com/idelchi/envprof)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/envprof)](https://goreportcard.com/report/github.com/idelchi/envprof)
[![Build Status](https://github.com/idelchi/envprof/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/envprof/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`envprof` is a CLI tool for managing named environment profiles in `YAML` or `TOML`.

Supports profile inheritance (layering) and importing of `.env` files.

## Features

- Define multiple environment profiles in a single YAML or TOML file, with inheritance and dotenv support
- List profiles, write to `.env` files or export to the current shell,
  execute a command or spawn a subshell with the selected environment

## Installation

For a quick installation, you can use the provided installation script:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/envprof/refs/heads/main/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
# list all profiles
envprof profiles
```

```sh
# list all variables in a profile with inheritance information
envprof --profile dev list -v
```

```sh
# list a specific variable
envprof --profile dev list HOST
```

```sh
# write profile to a file
envprof --profile dev write .env
```

```sh
# spawn a subshell with the environment loaded
envprof --profile dev shell
```

```sh
# export to current shell
eval "$(envprof --profile dev export)"
```

```sh
# Execute a command with the profile's environment
envprof --profile dev exec -- ls -la
```

## Format

Complex types (arrays, maps) are serialized as JSON; all other values are simple strings.

Non-scalars are emitted as minified JSON wrapped in single quotes.

### YAML

```yaml
dev:
  # Default profile to use when none is set with `--profile`
  default: true
  # Default output name for the `write` subcommand if not overridden by arguments
  output: development.env
  # Extend from other profiles
  extends:
    - staging
  env:
    HOST: localhost

staging:
  extends:
    - prod
  # Import dotenv files (relative to the current directory), supporting environment variables
  dotenv:
    - secrets.env
  env:
    HOST: staging.example.com
    DEBUG: true

prod:
  env:
    HOST: prod.example.com
    PORT: 80
    DEBUG: false
```

The `env` key alternatively accepts a sequence of key-value pairs:

```yaml
dev:
  env:
    - HOST=localhost
    - DEBUG=true
```

### TOML

```toml
[dev]
# Default profile to use when none is set with `--profile`
default = true
# Default output name for the `write` subcommand if not overridden by arguments
output = 'development.env'
# Extend from other profiles
extends = ['staging']
[dev.env]
HOST = 'localhost'

[staging]
extends = ['prod']
# Import dotenv files (relative to the current directory), supporting environment variables
dotenv = ['secrets.env']
[staging.env]
DEBUG = true
HOST = 'staging.example.com'

[prod.env]
DEBUG = false
HOST = 'prod.example.com'
PORT = 80
```

## Inheritance Behavior

Inheritance is resolved in order: later profiles override earlier ones.
Within each profile, `dotenv` files load before that profile’s environment variables.

As an example, running `envprof --profile dev write .env` with the previous YAML definition
as well as a sample `secrets.env`:

```sh
TOKEN=secret
```

produces the following `.env` file:

```sh
# Active profile: "dev"
DEBUG=true
HOST=localhost
PORT=80
TOKEN=secret
```

`envprof --profile dev list -v` shows the variables and their origins:

```sh
DEBUG=true              (inherited from "staging")
HOST=localhost
PORT=80                 (inherited from "prod")
TOKEN=secret            (inherited from "secrets.env")
```

The layering order here is:

```sh
prod -> secrets.env -> staging -> dev
```

from lowest to highest priority (left to right).

## Flags

All commands accept the following flags:

```sh
--file, -f      - Specify the profile file(s) to load
--profile, -p   - Specify the profile to use
--verbose, -v   - Increase verbosity
```

`--file` can be used to specify a file (or a list of fallback files) to load.
Defaults to the first found among `envprof.yaml`, `envprof.yml`, or `envprof.toml`, unless `ENVPROF_FILE` is set.

`--profile` specifies the profile to activate. If no profile is specified,
the [default profile](#yaml) will be used (if it exists).

`--verbose` increases verbosity, see subcommands for details.

## Subcommands

For details, run `envprof <command> --help` for the specific subcommand.

<details>
<summary><strong>profiles / profs</strong> — List all profiles</summary>

- **Usage:**
  - `envprof profiles`

- **Flags:**
  - `--verbose`, `-v` – Mark active profile with asterisk

</details>

<details>
<summary><strong>list / ls</strong> — List profile or the value of a variable in a profile</summary>

- **Usage:**
  - `envprof list [flags] [variable]`

- **Flags:**
  - `--oneline`, `-o` – Emit variables on a single line (implies `--verbose=false`)
  - `--verbose`, `-v` – Show variable origins

</details>

<details>
<summary><strong>export / x</strong> — Export profile to stdout</summary>

- **Usage:**
  - `envprof export [flags]`

- **Flags:**
  <!-- markdownlint-disable MD038 -->
  - `--prefix <string>` – String to prefix variables (default: `export `)
  <!-- markdownlint-enable MD038 -->

</details>

<details>
<summary><strong>write / w</strong> — Write profile(s) to file(s)</summary>

- **Usage:**
  - `envprof write [flags] [file]`

- **Flags:**
  - `--all`, `-a` – Write all profiles

</details>

<details>
<summary><strong>shell / sh</strong> — Spawn a subshell with profile</summary>

- **Usage:**
  - `envprof shell [flags]`

- **Flags:**
  - `--shell <shell>`, `-s <shell>` – Force shell (default empty string -> detected)
  - `--isolate`, `-i` – Prevent inheriting current shell variables
  - `--path`, `-p` – Include the current PATH in the environment

</details>

<details>
<summary><strong>exec / ex</strong> — Execute a command with profile</summary>

- **Usage:**
  - `envprof exec [flags] -- <command> [args...]`

- **Flags:**
  - `--isolate`, `-i` – Prevent inheriting current shell variables
  - `--path`, `-p` – Include the current PATH in the environment

</details>

## Shell integration

When using the `shell` subcommand, `envprof` sets `ENVPROF_ACTIVE_PROFILE` in the environment.

This variable is used to detect if you’re already in an `envprof` subshell, preventing nested sessions.

### Prompt

Use `ENVPROF_ACTIVE_PROFILE` to customize a `starship` prompt:

**`starship.toml`**

```toml
[env_var.envprof]
variable = "ENVPROF_ACTIVE_PROFILE"
format = '[\[envprof: $env_value\]]($style)'
style = 'bold bright-green'
```

## Demo

![Demo](assets/gifs/envprof.gif)
