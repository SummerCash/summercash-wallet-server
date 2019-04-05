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

#### Creating a New Account With an Existing Address (pseudo-code)

```Go
request := {
    "address": "your_existing_address",
} // Replace value with your SummerCash address

http.Post("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```

Responds with:

```JSON
{
    "name": "username",
    "password_hash": "387fud8739d7faef=",
    "address": "0x123456",
}
```

#### Creating a New Account

```Go
request := {} // Empty request so summercash-wallet-server generates a new address for us

http.Post("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```

Responds with:

```JSON
{
    "name": "username",
    "password_hash": "387fud8739d7faef=",
    "address": "0x123456",
}
```

#### Updating an Account's Password

```Go
request := {
    "old_password": "old_password", // Replace with the account's current password
    "new_password": "new_password", // Replace with the desired new password
}

http.Put("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```

Responds with:

```JSON
{
    "name": "username",
    "password_hash": "updated_hash",
    "address": "0x123456",
}
```

#### Fetching Account Details

```Go
request := {} // Empty request

http.Get("https://localhost:443/api/accounts/username", request) // Replace 'username' in '/username' with the desired username
```

Responds with:

```JSON
{
    "name": "username",
    "password_hash": "387fud8739d7faef=",
    "address": "0x123456",
}
```

### Transactions

#### Creating, Signing, and Publishing a New Transaction (pseudo-code)

```Go
request := {
    "username": "sender_username", // Replace with username of wallet to send from
    "password": "account_password", // Password of account to send from
    "recipient": "recipient_username_or_address", // Replace with recipient username or address
    "amount": 0, // Replace with amount to send w/tx
    "payload": "message_to_send_with_tx", // Replace w/transaction payload (e.g. contract call, message, etc...)
}
```

Responds with:

```JSON
{
    "nonce": 0,
    "sender": "0x123456",
    "recipient": "0x654321",
    "amount": 0,
    "payload": "d78fds=",
    "signature": {
        "SerializedPublicKey": "DYfs87v997awe...",
        "V": "7bPJcyxV1VIQGWO/kavPpYUG66mP7n1qnn2fRnV6pBk=",
        "R": 012345678910,
        "S": 012345678910,
    },
    "time": "2019-04-04T22:22:03.084703Z",
    "contract": null,
    "is-init-contract": false,
    "genesis": false,
    "logs": null,
    "hash": "0x123456",
}
```
