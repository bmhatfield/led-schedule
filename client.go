package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const fade = `brightness 1,%d
render
`

type Client struct {
	Address        string
	FullBrightness int
	FadeInterval   time.Duration
	IsOn           bool
}

func (c Client) conn() (net.Conn, error) {
	return net.Dial("tcp", c.Address)
}

func (c Client) setColor(color string, conn net.Conn) error {
	_, err := conn.Write([]byte(fmt.Sprintf("fill 1,%s\n", color)))
	return err
}

func (c *Client) FadeUp() error {
	log.Println("connecting...")
	conn, err := c.conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := c.setColor("BB5511", conn); err != nil {
		return err
	}

	log.Println("fading up...")
	if _, err := conn.Write([]byte(
		fmt.Sprintf("fade 1,0,%d,%d,1\n", c.FullBrightness, c.FadeInterval.Milliseconds()),
	)); err != nil {
		return err
	}

	c.IsOn = true
	return nil
}

func (c Client) FadeDown() error {
	log.Println("connecting...")
	conn, err := c.conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("fading down...")
	if _, err = conn.Write([]byte(
		fmt.Sprintf("fade 1,%d,0,%d,1\n", c.FullBrightness, c.FadeInterval.Milliseconds()),
	)); err != nil {
		return err
	}

	c.IsOn = false
	return nil
}
