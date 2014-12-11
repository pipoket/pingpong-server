package main


import (
    "log"
    "net"
)


type udpServer struct {
    bindIp string
    bindPort int
    conn *net.UDPConn
}


func (server *udpServer) Start() {
    log.Printf("UDPServer::start")

    addr := net.UDPAddr{
        IP: net.ParseIP(server.bindIp),
        Port: server.bindPort,
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        panic(err)
    }
    server.conn = conn
}


func (server *udpServer) ServeForever() {
    for {
        var buffer [1024]byte
        rlen, raddr, err := server.conn.ReadFrom(buffer[:])
        if err != nil {
            log.Printf("ERROR while reading")
            return
        }
        server.HandlePacket(raddr, rlen, buffer[:])
    }
}


func (server *udpServer) HandlePacket(raddr net.Addr, rlen int, buffer []byte) {
    go func(raddr net.Addr, rlen int, buffer []byte) {
        log.Printf("Read %d bytes from %s", rlen, raddr.String())
    }(raddr, rlen, buffer)
}


func main() {
    log.Printf("Hello")
    server := udpServer{bindIp: "0.0.0.0", bindPort: 7000}

    server.Start()
    server.ServeForever()
    log.Printf("Bye")
}
