# tcp-keepalive
tcp keepalive measuring tool

## install

```
git clone https://github.com/Openresty-App/tcp-keepalive
cd tcp-keepalive

linux:
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tcp-keepalive

macos:
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o tcp-keepalive

freebsd:
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o tcp-keepalive

```

## Sample
```
./tcp-keepalive -h

./tcp-keepalive --host 127.0.0.1 --port 9527
```
