package subcmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tcp-keepalive/utils"
	"time"

	"github.com/felixge/tcpkeepalive"
	"github.com/urfave/cli"
)

func heartbeat(host string, port, idleTime, count, interval int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	idleTimeDuration := time.Duration(time.Duration(idleTime) * time.Second)
	intervalDuration := time.Duration(time.Duration(interval) * time.Second)
	log.Printf("idleTime:%+v, interval:%+v, count:%d", idleTimeDuration, intervalDuration, count)

	if err := tcpkeepalive.SetKeepAlive(conn, idleTimeDuration, count, intervalDuration); err != nil {
		return err
	}

	for {
		if utils.IsClosed(conn) {
			return errors.New("conn reset")
		}
		time.Sleep(time.Second)
	}

	return nil
}

func doHeartbeat(host string, port, idleTime, count, interval int, ch chan int) {
	if err := heartbeat(host, port, idleTime, count, interval); err != nil {
		log.Printf("heartbeat err :%s", err.Error())
		time.Sleep(time.Second)
		ch <- 1
	}
}

func doTcpDump(device string, filter string, ch chan int, cap chan Capture) {
	tcpdump(device, filter, cap)
}

func th(c *cli.Context) {
	host := c.String("host")
	port := c.Int("port")
	device := c.String("interface")
	idleTime := c.Int("idleTime")
	count := c.Int("count")
	interval := c.Int("interval")

	ch := make(chan int)
	cap := make(chan Capture)
	met := utils.NewMetric()

	filter := fmt.Sprintf("tcp and port %d", port)
	go doHeartbeat(host, port, idleTime, count, interval, ch)

	time.Sleep(time.Duration(idleTime) * time.Second)
	go doTcpDump(device, filter, ch, cap)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)

	var a, b time.Time
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
					a = tcp.GetTime()
					getOne = true
				} else {
					if getOne {
						b = tcp.GetTime()
						met.Put(b.Sub(a).Seconds())
					} else {
						a = time.Time{}
					}

					getOne = false
				}
			}

		case c := <-ch:
			if c == 1 {
				log.Printf("heartbeat restart ...")
				time.Sleep(time.Second * 3)
				go doHeartbeat(host, port, idleTime, count, interval, ch)
			} else if c == -1 {
				log.Printf("get error")

				break L
			}
		}
	}

	log.Printf("max:%.6f, mean:%.6f, min:%.6f", met.Max(), met.Mean(), met.Min())
	return
}

func Th() cli.Command {
	return cli.Command{
		Name:  "th",
		Usage: "tcp heart and statistic",
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
			cli.IntFlag{
				Name:  "interval",
				Value: 3,
				Usage: "interval",
			},
		},
		Action: th,
	}
}
