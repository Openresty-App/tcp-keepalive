package subcmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

func probe(c *cli.Context) {
	host := c.String("host")
	port := c.Int("port")
	device := c.String("interface")
	idleTime := c.Int("idleTime")
	retranCount := c.Int("retranCount")
	interval := 3

	ch := make(chan int)
	cap := make(chan Capture)

	filter := fmt.Sprintf("host %s and tcp and port  %d", host, port)
	go doHeartbeat(host, port, idleTime, retranCount, interval, ch)

	time.Sleep(time.Duration(idleTime) * time.Second)
	go doTcpDump(device, filter, ch, cap)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)

	res := "fail"
	getOne := false

L:
	for {
		select {
		case s := <-sig:
			log.Println("退出信号", s)
			break L

		case d := <-cap:
			if tcp, ok := d.(*CaptureTcpData); ok {
				fmt.Printf(tcp.ToString())

				// req
				if int(tcp.GetDport()) == port {
					getOne = true
				} else {
					if getOne && int(tcp.GetSport()) == port {
						res = "ok"
					}
					break L
				}
			}

		case c := <-ch:
			if c == 1 {
				break L
			} else if c == -1 {
				fmt.Printf("get error")
				break L
			}
		}
	}
	fmt.Printf("tcp heartbeat %s", res)
	return
}

func Probe() cli.Command {
	return cli.Command{
		Name:  "probe",
		Usage: "Check whether the server supports heartbeat",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "interface,i",
				Value: "eth0",
				Usage: "interface",
			},
			cli.StringFlag{
				Name:  "host,H",
				Value: "127.0.0.1",
				Usage: "host",
			},
			cli.IntFlag{
				Name:  "port, P",
				Value: 80,
				Usage: "port",
			},
			cli.IntFlag{
				Name:  "idleTime",
				Value: 3,
				Usage: "idleTime",
			},
			cli.IntFlag{
				Name:  "retranCount",
				Value: 10,
				Usage: "retranCount",
			},
		},
		Action: probe,
	}
}
