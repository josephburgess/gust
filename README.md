# gust

A simple cli weather tool built with Go that provides weather forecasts in the terminal.

<p align="center">
<img src="https://github.com/user-attachments/assets/d5c8e7cc-c43c-4263-a516-89f01bb26a24" width="600">
</p>

## Installation

### Option 1: Homebrew (macOS)

The easiest way to install gust at the moment is with Homebrew:

```sh
brew tap josephburgess/formulae
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

## Configuration Flags

_These flags modify user config / settings and don't display weather by themselves_

| Short | Long               | Description                                                |
| ----- | ------------------ | ---------------------------------------------------------- |
| `-h`  | `--help`           | Show help                                                  |
| `-S`  | `--setup`          | Run the setup wizard                                       |
| `-A`  | `--api=STRING`     | Set custom API server URL (mostly for development)         |
| `-C`  | `--city=STRING`    | Specify city name                                          |
| `-D`  | `--default=STRING` | Set a new default city                                     |
| `-U`  | `--units=STRING`   | Set default temperature units (metric, imperial, standard) |
| `-L`  | `--login`          | Authenticate with GitHub                                   |
| `-K`  | `--api-key`        | Set your api key (either gust, or openweathermap)          |

## Display Flags

_These flags control how weather information is displayed_

| Short | Long         | Description                                   |
| ----- | ------------ | --------------------------------------------- |
| `-a`  | `--alerts`   | Show weather alerts                           |
| `-c`  | `--compact`  | Show today's compact weather view             |
| `-d`  | `--detailed` | Show today's detailed weather view            |
| `-f`  | `--full`     | Show today, 5-day and weather alert forecasts |
| `-r`  | `--hourly`   | Show 24-hour (hourly) forecast                |
| `-y`  | `--daily`    | Show 5-day forecast                           |

## Authentication

gust uses a proxy api I set up and host privately, [breeze](http://github.com/josephburgess/breeze), to fetch weather data. This keeps the setup flow pretty frictionless for new users.

When you first run gust:

1. A setup wizard will guide you through the initial configuration
2. You'll be prompted to choose an authentication method:
   - GitHub OAuth (one click sign up/in)
   - Your own OpenWeather Map API key if you prefer not to use Oauth or need much higher rate limits
     - User submitted keys will need to be eligible for the [One Call API 3.0](https://openweathermap.org/api/one-call-3#how)
     - The first 1000 calls every day are free but they ask for CC info to get a key
3. If you choose GitHub OAuth:
   - Your default browser will open to complete authentication
   - No need to manually obtain or manage API keys
4. Your credentials will be securely stored locally for future use in `~/.config/gust/auth.json`

After this one-time setup, authentication happens automatically whenever you use the app.

## Troubleshooting

If you encounter any auth issues, you can re-run the setup wizard or use the `-L / --login` (Oauth) `-K / --api-key` (api key) flags to re-set your key or check the local config files.
