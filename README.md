# **Gust - Weather cli in Go**

Gust is a simple command-line weather tool I'm using to learn/play with Golang.

## **Installation**

### **1. Clone the repository**

```sh
git clone https://github.com/josephburgess/gust.git
cd gust
```

### **2. Set up environment variables**

You'll need an OpenWeather API key. Get one [here](https://home.openweathermap.org/api_keys).

Create a `.env` file in the root directory:

```sh
echo "OPENWEATHER_API_KEY=your_api_key_here" > .env
```

### **3. Install dependencies**

```sh
go mod tidy
```

### **4. Build the binary**

```sh
go build -o gust ./cmd/gust
```

### **5. Make it globally accessible**

Move the binary to a directory in your `$PATH`:

```sh
mv gust /usr/local/bin/
```

## **Usage**

Run the command with a city name:

```sh
gust london
```

### **Example Output**
![Screenshot 2025-02-26 at 00 10 45@2x](https://github.com/user-attachments/assets/90ce75c0-6cba-40fa-9aae-16cd53064f52)

## **Testing**

Run tests using:

```sh
go test ./...
```

## **License**

MIT License. Free to use and modify.
