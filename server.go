package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var Map map[string]string

//var mut sync.Mutex
var mut sync.RWMutex

func main() {
	if len(os.Args) != 2 {
		fmt.Println("format:-server ipaddr:port")
		return
	}
	//ch := make(chan bool)
	Map = make(map[string]string)

	socket := os.Args[1]
	addr_port, err := net.ResolveTCPAddr("tcp", socket)
	if err != nil {
		fmt.Println("Address format is ipaddr:port_num")
	}
	ln, err := net.ListenTCP("tcp", addr_port)
	if err != nil {
		fmt.Println("use different port number, given port already in use", err.Error())
	}

	fmt.Println("Listening at addr:port", addr_port)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print(err)
			return
		}
		//wg.Add(1)
		//go echo_srv(conn, wg)
		go echo_srv(conn)

		//wg.Wait()
	}
}

func Get_Val(key string) string {
	return Map[key]
}

func Set_Val(key string, val string) {
	Map[key] = val
	return
}

func Update_Val(key string, val string) {
	Map[key] = val
	return
}

func Delete(key string) {
	delete(Map, key)
	return
}

func Rename(key1 string, key2 string) bool {
	ok := true
	_, ok = Map[key2] //first check key2 must not consist any-value.
	if ok {
		return !ok
	}
	Map[key2] = Map[key1]
	delete(Map, key1)
	return !ok
}

func echo_srv(c net.Conn) {
	defer c.Close()

	for {
		msg := make([]byte, 1000)

		n, err := c.Read(msg)
		if err == io.EOF {
			fmt.Printf("SERVER: received EOF (%d bytes ignored)\n", n)
			return
		} else if err != nil {
			fmt.Printf("ERROR: read\n")
			fmt.Print(err)
			return
		}
		fmt.Printf("SERVER: received %v bytes\n", n)

		//to_send :=string(msg);
		str := strings.Fields(string(msg))

		fmt.Printf("msg=%s\n", msg)
		//fmt.Println(str[0]=="get",str[1])
		action := str[0]
		key := str[1]
		switch {

		case action == "set":
			val := str[2]
			mut.Lock()
			//Map[key] = val
			Set_Val(key, val)
			mut.Unlock()
			fmt.Printf("m[%s]=%s", str[1], str[2])

		case action == "update":
			val := str[2]
			mut.Lock()
			//Map[key] = val
			Update_Val(key, val)
			mut.Unlock()
			fmt.Printf("m[%s]=%s", str[1], str[2])

		case action == "get":
			fmt.Printf("s Map[%s]=%s e\n", str[1], Map[key])
			fmt.Printf("s Map[%s]=%s e\n", str[1], Map["abc"])
			//temp:=string(strings.TrimSpace(str[1]));
			mut.RLock()
			//temp := []byte(Map[key])
			temp := []byte(Get_Val(key))
			mut.RUnlock()
			fmt.Println(temp)
			n, err = c.Write(temp)
			if n == 0 {
				n, err = c.Write([]byte("\n"))
			}

			if err != nil {
				fmt.Printf("ERROR: write\n")
				fmt.Print(err)
				return
			}
			fmt.Printf("SERVER: sent %v bytes\n", n)
		}
	}
}
