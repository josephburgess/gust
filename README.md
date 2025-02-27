# gust
Gust is a simple command-line weather tool I'm using to learn/play with Golang.
<p align="center">
<img src="https://github.com/user-attachments/assets/9aa70f5e-cfed-4265-a6ee-54f4e5f59ea8" width="400">
</p>
## **Installation**

### **1. Clone the repository**

```sh
git clone https://github.com/josephburgess/gust.git
cd gust
```

### **2. Set up .env**

gust uses openweather, you can get a key [here](https://home.openweathermap.org/api_keys), then create a `.env` with your key:
```sh
echo "OPENWEATHER_API_KEY=your_api_key_here" > .env
```

### **3. Install deps**

```sh
go mod tidy
```

### **4. Build the binary**

```sh
go build -o gust ./cmd/gust
```

### **5. Install**

Move the binary to a dir in your `$PATH`:

```sh
mv gust /usr/local/bin/
```

## **Usage**

Run the command with a city name:

```sh
gust london
```

## **Testing**

Run tests using:

```sh
go test ./...
```

## **License**

MIT License. Free to use and modify.
