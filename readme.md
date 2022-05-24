# Getting Start

**Do not push the binary file.**

## Mac Os

Install go v1.18

```bash
$ brew install go
```

Install dependency

```bash
$ go mod tidy
```

Run demo server

```bash
$ go run main.go
```

## Linux

Build binary for GNU/Linux x64 on your Mac OS.

```bash
$ GOOS=linux GOARCH=amd64 go build main.go
```

Scp the binary to the server.

```bash
$ scp ./main root@ipAdreess:/root/app
```

Run the main on the **server side**.

```bash
$ ./main
```

