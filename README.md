# SummerCash Wallet Server

The golang react-summercash webserver and backend.

## Installation

```zsh
go install github.com/SummerCash/summercash-wallet-server
```

## Usage

Starting the webserver:

```zsh
summercash-wallet-server
```

or

```zsh
go run main.go
```

### Serving Content

To serve static HTML content with summercash-wallet-server, simply copy all necessary content into a content/ folder in the summercash-wallet-server root (or specify via the --content-dir flag).

## APIs

| URI                                      | Name             | Description                                                                          |
| ---------------------------------------- | ---------------- | ------------------------------------------------------------------------------------ |
| <https://localhost:443/api/accounts>     | Accounts API     | An API for creating, managing, and fetching SummerCash account details.              |
| <https://localhost:443/api/transactions> | Transactions API | An API for creating, signing, and publishing transactions on the SummerCash network. |
