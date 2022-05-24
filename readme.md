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
$ go run demo.go
```

## Linux

Build binary for GNU/Linux x64 on your Mac OS.

```bash
$ GOOS=linux GOARCH=amd64 go build demo.go
```

Scp the binary to the server.

```bash
$ scp ./demo root@ipAdreess:/root/app
```

Run the demo on the **server side**.

```bash
$ ./demo
```

