package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverip string, serverport int) *Client {
	client := &Client{
		ServerIp:   serverip,
		ServerPort: serverport,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverip, serverport))
	if err != nil {
		fmt.Println("net.dial error:", err)
	}

	client.conn = conn

	return client
}

func (c *Client) DealResponse() {
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更换名字")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的数字")
		return false
	}
}

func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容，exit退出<<<<<")
		fmt.Scanln(&chatMsg)
	}
}

func (c *Client) QueryUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *Client) PrivateChat() {

	var remoteName string
	var chatMsg string

	c.QueryUsers()
	fmt.Println(">>>>>请输入聊天对象（用户名），exit退出<<<<<")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>请输入消息内容，exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("write error:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>请输入消息内容，exit退出")
			fmt.Scanln(&chatMsg)
		}

		remoteName = ""
		fmt.Println(">>>>>请输入聊天对象（用户名），exit退出<<<<<")
		fmt.Scanln(&remoteName)
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名<<<<<")
	_, _ = fmt.Scanln(c.Name)

	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("write error :", err)
		return false
	}

	return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}

		switch c.flag {
		case 1:
			c.PublicChat()
			break
		case 2:
			c.PrivateChat()
			break
		case 3:
			c.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口号")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("failed to connect servers...")
		return
	}

	go client.DealResponse()

	fmt.Println("success to connect servers...")

	client.Run()
}
