## pScan: A CLI port scanner tool built in Golang

A Detailed documentation about the tool and its usage is available [here](docs/pScan.md)

## How to run

### Prerequisites
- Ensure you have [Golang](https://go.dev/doc/install) installed on your machine.

### Installation
1. Clone the repository
```
git clone https://github.com/auxence-m/port-scanner.git
```

2. Navigate to the Project Directory
```
cd port-scanner
```

3. Install Dependencies
```
go mod tidy
```

4. Build the Application
```
go build
``` 

After building, you'll find the `pScan` executable (`pScan.exe` on Windows) in your project directory.

5. Run the Application
```
pScan [command] --flag
```