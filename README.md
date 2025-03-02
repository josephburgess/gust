# gust

A simple cli weather tool built with Go that provides weather forecasts in the terminal.

<p align="center">
<img src="https://github.com/user-attachments/assets/76695b8d-5e37-45a3-89cd-2d5b3401c323" width="600">
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

Simply run the command with a city name:

```sh
gust london
```

### Optional Parameters

```
gust --city london                   # Specify the name of a city other than your default
gust --default london                # Set a new default city
gust --login                         # Authenticate with GitHub
gust --logout                        # Log out and remove authentication
gust --api https://custom.api.server # Set custom API server URL if you're using breeze locally
gust --daily london                  # Show daily forecast only
gust --hourly london                 # Show hourly forecast only
gust --alerts london                 # Show weather alerts only
gust --full london                   # Show full weather report including daily and hourly forecasts (might be a lot to read!)
```

## Authentication

Gust uses a proxy api, [breeze](http://github.com/josephburgess/breeze), to fetch weather data. This means you no longer need your own OpenWeather API key and can just authenticate with SSO using GitHub.

When you first run Gust, it will:

1. Prompt you to authenticate with GitHub
2. Open your browser for the authentication process
3. Automatically store your API key for future use

After the initial setup, authentication is fully automatic and your API key will be securely stored locally.
