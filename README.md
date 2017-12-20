# tcp-keepalive
tcp keepalive probe or measure tool

## Depends
* Linux
  * [libpcap](http://www.tcpdump.org/)
* Windows
  * [winpcap](https://www.winpcap.org/)

## Install

```lang=shell
go get github.com/felixge/tcpkeepalive
go get github.com/urfave/cli

git clone https://github.com/Openresty-App/tcp-keepalive
cd tcp-keepalive


go build -o tcp-keepalive
```

## Measure
```
#./tcp-keepalive help measure
NAME:
   tcp-keepalive measure - Send tcp heart and statistic

USAGE:
   tcp-keepalive measure [command options] [arguments...]

OPTIONS:
   --interface value, -i value  interface (default: "eth0")
   --host value, -H value       host (default: "127.0.0.1")
   --port value, -P value       port (default: 80)
   --idleTime value             idleTime (default: 3)
   --retranCount value          retranCount (default: 10)
   --interval value             interval (default: 3)

./tcp-keepalive measure -H 172.28.32.220 -P 8888 -i en5
```

## Probe
```
#./tcp-keepalive help probe
NAME:
   tcp-keepalive probe - Check whether the server supports heartbeat

USAGE:
   tcp-keepalive probe [command options] [arguments...]

OPTIONS:
   --interface value, -i value  interface (default: "eth0")
   --host value, -H value       host (default: "127.0.0.1")
   --port value, -P value       port (default: 80)
   --idleTime value             idleTime (default: 3)
   --retranCount value          retranCount (default: 10)

./tcp-keepalive probe -H 172.28.32.220 -P 8888 -i en5
```

## Team

This tool is maintained by the following person(s) and a bunch of [awesome contributors](https://github.com/Openresty-App/tcp-keepalive/graphs/contributors).

## License

[MIT License](./LICENSE)


