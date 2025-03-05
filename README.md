# gust

A simple cli weather tool built with Go that provides weather forecasts in the terminal.

<p align="center">
<img src="https://github.com/user-attachments/assets/d5c8e7cc-c43c-4263-a516-89f01bb26a24" width="600">
</p>

## Installation

### Option 1: Homebrew (macOS)

The easiest way to install gust at the moment is with Homebrew:

```sh
brew tap josephburgess/tools
brew install gust
```

### Option 2: Manual Installation

#### 1. Clone the repository

```sh
git clone https://github.com/josephburgess/gust.git
cd gust
```

#### 2. Install dependencies

```sh
go mod tidy
```

#### 3. Build the binary

```sh
go build -o gust ./cmd/gust
```

#### 4. Install

Move the binary to a directory in your `$PATH`:

```sh
mv gust /usr/local/bin/
```

## Usage

### Basic Commands

```bash
# Get weather for your default city
gust

# Get weather for a specific city
gust london
```

<p align="center">
  <img src="https://github.com/user-attachments/assets/76695b8d-5e37-45a3-89cd-2d5b3401c323" width="600">
</p>

### Configuration Flags

_These flags modify settings and don't display weather by themselves_

| Short | Long               | Description                                                |
| ----- | ------------------ | ---------------------------------------------------------- |
| `-h`  | `--help`           | Show help                                                  |
| `-s`  | `--setup`          | Run the setup wizard                                       |
| `-D`  | `--default=STRING` | Set a new default city                                     |
| `-u`  | `--units=STRING`   | Set default temperature units (metric, imperial, standard) |
| `-l`  | `--login`          | Authenticate with GitHub                                   |
| `-A`  | `--api=STRING`     | Set custom API server URL (mostly for development)         |

### Display Flags

_These flags control how weather information is displayed_

| Short | Long         | Description                                   |
| ----- | ------------ | --------------------------------------------- |
| `-c`  | `--compact`  | Show today's compact weather view             |
| `-d`  | `--detailed` | Show today's detailed weather view            |
| `-y`  | `--daily`    | Show 5-day forecast                           |
| `-h`  | `--hourly`   | Show 24-hour (hourly) forecast                |
| `-a`  | `--alerts`   | Show weather alerts                           |
| `-f`  | `--full`     | Show today, 5-day and weather alert forecasts |

## Authentication

Gust uses a proxy api I set up and host privately, [breeze](http://github.com/josephburgess/breeze), to fetch weather data. This keeps the setup flow pretty frictionless for new users.

When you first run Gust, it will:

1. Take you through a setup wizard
2. Prompt you to authenticate with GitHub
3. Open your browser for the authentication process
4. Automatically store your API key for future use

After the initial setup, authentication is fully automatic and your API key will be securely stored locally.

If for some reason any of this fails you should be able to attempt to setup/login again with the flags above.
