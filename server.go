package server

import (
        "fmt"
        "net"
	"io"
	//"sync"
	"strings"
)

var m  map[string]string

func main() {
	
    m = make(map[string]string)
    //var wg sync.WaitGroup

        service := "0.0.0.0:12345"
        tcpAddr, err := net.ResolveTCPAddr("tcp", service)
        checkError(err)
        ln, err := net.ListenTCP("tcp", tcpAddr)
        checkError(err)
        fmt.Println("Server listening at port number ",tcpAddr)
    if err != nil {
            fmt.Print(err)
            return
    }
    defer ln.Close()
for{
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

func checkError(err error) {
        if err != nil {
                fmt.Println("Fatal error ", err.Error())
        }
}

//func echo_srv(c net.Conn, wg sync.WaitGroup) {
func echo_srv(c net.Conn) {
    defer c.Close()

    for {
        msg := make([]byte, 1000)

        n, err := c.Read(msg)
        if err == io.EOF {
            fmt.Printf("SERVER: received EOF (%d bytes ignored)\n", n)
            return
        } else  if err != nil {
            fmt.Printf("ERROR: read\n")
            fmt.Print(err)
            return
        }
        fmt.Printf("SERVER: received %v bytes\n", n)

	//to_send :=string(msg);
	str := strings.Fields(string(msg));

	fmt.Printf("msg=%s\n",msg);
	//fmt.Println(str[0]=="get",str[1])
        key := str[1]
	if str[0]=="set"{
		val := str[2]
		m[key]=val
		fmt.Printf("m[%s]=%s",str[1],str[2])
	}
	if str[0]=="get" {
		fmt.Printf("s m[%s]=%s e\n",str[1],m[key])
		fmt.Printf("s m[%s]=%s e\n",str[1],m["abc"])
		//temp:=string(strings.TrimSpace(str[1]));
		temp:=[]byte(m[key]);
		fmt.Println(temp);
		//fmt.Println(m,"checking",m["h"],"m[1]",str[1],"m[temp]=",temp=="hel")
	        //n, err = c.Write([]byte(m[key]))
	        n, err = c.Write(temp)
                if n==0{
	                n,err = c.Write([]byte("\n"));
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

