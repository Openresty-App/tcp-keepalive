package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/felixge/tcpkeepalive"
	"github.com/urfave/cli"
)

func isClosed(c *net.TCPConn) bool {
	one := make([]byte, 1)
	c.SetReadDeadline(time.Time{})
	_, err := c.Read(one)
	if err == io.EOF {
		c.Close()
		c = nil
		return true
	} else {
		var zero time.Time
		c.SetReadDeadline(zero)

		return false
	}
}

func do(host string, port, idleTime, count, interval int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	log.Printf("connecting")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	log.Printf("connected")
	defer conn.Close()

	idleTimeDuration := time.Duration(time.Duration(idleTime) * time.Second)
	intervalDuration := time.Duration(time.Duration(interval) * time.Second)
	log.Printf("idleTime:%+v, interval:%+v, count:%d", idleTimeDuration, intervalDuration, count)

	if err := tcpkeepalive.SetKeepAlive(conn, idleTimeDuration, count, intervalDuration); err != nil {
		return err
	}

	for {
		if isClosed(conn) {
			return errors.New("conn reset")
		}
		time.Sleep(time.Second)
	}
	//
	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//s := <-c
	//fmt.Println("退出信号", s)
	return nil
}

func main() {
	app := cli.NewApp()
	var host string
	var port, idleTime, count, interval int

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Value:       "127.0.0.1",
			Usage:       "host",
			Destination: &host,
		},
		cli.IntFlag{
			Name:        "port",
			Value:       80,
			Usage:       "port",
			Destination: &port,
		},
		cli.IntFlag{
			Name:        "idleTime",
			Value:       1,
			Usage:       "idleTime",
			Destination: &idleTime,
		},
		cli.IntFlag{
			Name:        "count",
			Value:       10,
			Usage:       "count",
			Destination: &count,
		},
		cli.IntFlag{
			Name:        "interval",
			Value:       1,
			Usage:       "interval",
			Destination: &interval,
		},
	}

	app.Action = func(c *cli.Context) error {
		err := do(host, port, idleTime, count, interval)
		if err != nil {
			log.Println(err)
		}

		return nil
	}

	app.Run(os.Args)
}
