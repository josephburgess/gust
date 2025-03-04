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

Simply run the command:
<p align="center">
  <img src="https://github.com/user-attachments/assets/76695b8d-5e37-45a3-89cd-2d5b3401c323" width="600">
</p>

Or if you want to check a city other than the default you set, then just add the name:

```bash
gust london
```

### Optional Parameters

```
gust --city london                   # Specify the name of a city other than your default
gust --setup                         # Run the setup wizard again
gust --default london                # Set a new default city
gust --login                         # Authenticate with GitHub
gust --logout                        # Log out and remove authentication
gust --api https://...               # Set custom API server URL if you're using breeze locally
gust --api https://...               # Set custom API server URL if you're using breeze locally
gust --detailed                      # Show detailed forcast for today
gust --compact                       # Show compact forecast for today
gust --daily                         # Show 5 day forecast
gust --hourly                        # Show 24hr forecast
gust --alerts                        # Show weather alerts only
gust --full london                   # Show detailed weather, plus any alerts, plus 5 day forecast
```

## Authentication

Gust uses a proxy api I set up and host privately, [breeze](http://github.com/josephburgess/breeze), to fetch weather data. This keeps the setup flow pretty frictionless for new users.

When you first run Gust, it will:

1. Take you through a setup wizard
2. Prompt you to authenticate with GitHub
3. Open your browser for the authentication process
4. Automatically store your API key for future use

After the initial setup, authentication is fully automatic and your API key will be securely stored locally.

If for some reason any of this fails you should be able to attempt to setup/login again with the flags above.
