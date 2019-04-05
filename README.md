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

To serve static content with summercash-wallet-server, simply copy all necessary content into a content/ folder in the summercash-wallet-server root (or specify via the --content-dir flag).

## APIs

| URI                                      | Name             | Description                                                                          |
| ---------------------------------------- | ---------------- | ------------------------------------------------------------------------------------ |
| <https://localhost:443/api/accounts>     | Accounts API     | An API for creating, managing, and fetching SummerCash account details.              |
| <https://localhost:443/api/transactions> | Transactions API | An API for creating, signing, and publishing transactions on the SummerCash network. |

### Accounts

Creating a New Account With an Existing Address (pseudo-code):

```Go
request := {
    "address": "your_existing_address",
} // Replace value with your SummerCash address

http.Post("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```

Creating a New Account:

```Go
request := {} // Empty request so summercash-wallet-server generates a new address for us

http.Post("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```

Updating an Account's Password:

```Go
request := {
    "old_password": "old_password", // Replace with the account's current password
    "new_password": "new_password", // Replace with the desired new password
}

http.Put("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```

Fetching Account Details:

```Go
request := {} // Empty request

http.Get("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```
