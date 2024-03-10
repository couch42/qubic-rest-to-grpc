<br /><br />

qubic-rest-to-grcp

<br /><br />

## Building the app binary:
```bash
$ go build -o "server" "."
```

## Running the server:
```bash
$ ./server
2024/03/10 00:08:17 Server listening on port 7070
```

## Docker usage:
```bash
$ docker build -t ghcr.io/couch42/qubic-rest-to-grcp:latest .
$ docker run -p 7070:7070 ghcr.io/couch42/qubic-rest-to-grcp:latest
```
