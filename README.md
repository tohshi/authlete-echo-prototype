# Authlete-Echo-Prototype

Prototype of an authorization server using Echo and Authlete.  
This implementation is not perfect and should be kept for reference only.

## Getting Started

1. Install dependencies

```sh
$ go mod tidy
```

2. Set up environment variables

```sh
$ vi .env
```

3. Start server on `http://localhost:1323`

```sh
$ go run *.go
```

### Endpoints

| Endpoint               | Path                  |
| :--------------------- | :-------------------- |
| Authorization Endpoint | `/auth/authorization` |
| Token Endpoint         | `/auth/token`         |

### Users

| Login ID | Password  | Consent Required |
| :------- | :-------- | :--------------- |
| `user1`  | `passwd2` | `true`           |
| `user2`  | `passwd2` | `false`          |
