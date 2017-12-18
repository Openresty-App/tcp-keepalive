# tcp-keepalive
tcp keepalive measuring tool

## Install

```lang=shell
go get github.com/felixge/tcpkeepalive
go get github.com/urfave/cli

git clone https://github.com/Openresty-App/tcp-keepalive
cd tcp-keepalive

linux:
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tcp-keepalive

macos:
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o tcp-keepalive

freebsd:
CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o tcp-keepalive

```

## Example
```
[root@localhost tcp-keepalive]# ./tcp-keepalive -h
NAME:
   tcp-keepalive - A new cli application

USAGE:
   tcp-keepalive [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value      host (default: "127.0.0.1")
   --port value      port (default: 80)
   --idleTime value  idleTime (default: 3)
   --count value     count (default: 10)
   --interval value  interval (default: 1)
   --help, -h        show help
   --version, -v     print the version

[root@localhost tcp-keepalive]# ./tcp-keepalive --host www.baidu.com --port 80 --idleTime 10
2017/12/09 14:50:53 connecting
2017/12/09 14:50:53 connected
2017/12/09 14:50:53 idleTime:10s, interval:1s, count:10
2017/12/09 14:51:14 conn reset
```
