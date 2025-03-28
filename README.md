# kryptos

A command-line tool for managing encrypted keys stored in Redis, written in Go. It provides functionality to delete, list, and get keys securely. This tool was made for learning purposes and 

## Installation

To install `kryptos`, use the following steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/chxmxii/kryptos.git
    cd kryptos
    ```

2. Install by running:
    ```bash
    make install
    ```

## Usage

`kryptos` provides several commands to manage encrypted secrets in Redis. Below are the available commands:

### List keys
To list all stored secrets:
```bash
kryptos list
```

### Get a key
To retrieve a specific key:
```bash
kryptos get <key> -i <dbidx> -k <decryption_key>
```

### Delete a key
To delete a key:
```bash
kryptos del <key>
```

### Add a new key
To add a new key:
```bash
kryptos put <key>:<value>  -i <dbidx> -k <encryption_key>
```

### Help
For detailed usage of any command, use the `--help` flag:
```bash
kryptos <command> --help
```

Make sure to configure your Redis/KeyDB connection settings before using the tool. You can do this by setting the appropriate environment variables or updating the configuration file.