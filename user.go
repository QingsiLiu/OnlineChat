package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	useraddr := conn.RemoteAddr().String()

	user := &User{
		Name: useraddr,
		Addr: useraddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	go user.ListenMessage()

	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		if _, err := u.conn.Write([]byte(msg + "\n")); err != nil {
			fmt.Println("write msg error:", err)
		}
	}
}

func (u *User) Online() {
	u.server.maplock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.maplock.Unlock()

	u.server.BroadCast(u, "is online")
}

func (u *User) Offline() {
	u.server.maplock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.maplock.Unlock()

	u.server.BroadCast(u, "is offline")
}

func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.maplock.Lock()
		for _, user := range u.server.OnlineMap {
			onlinemsg := "[" + user.Addr + "]" + user.Name + ":" + "is online...\n"
			u.SendMsg(onlinemsg)
		}
		u.server.maplock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("new name is used")
		} else if len(msg) > 4 && msg[:3] == "to|" {
			remoteName := strings.Split(msg, "|")[1]
			if remoteName == "" {
				u.SendMsg("error")
				return
			}

			remoteUser, ok := u.server.OnlineMap[remoteName]
			if !ok {
				u.SendMsg("not exists")
				return
			}

			content := strings.Split(msg, "|")[2]
			if content == "" {
				u.SendMsg("none msg")
				return
			}
			remoteUser.SendMsg(u.Name + "to you:" + content)
		} else {
			u.server.maplock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.maplock.Unlock()

			u.Name = newName
			u.SendMsg("update username success:" + u.Name + "\n")
		}
	} else {
		u.server.BroadCast(u, msg)
	}
}

func (u *User) SendMsg(msg string) {
	if _, err := u.conn.Write([]byte(msg)); err != nil {
		fmt.Println("write msg error:", err)
	}
}
