package pkg

import (
	"io"
	"log"
	"net"

	"github.com/askolesov/obfsproxy/pkg/codec"
)

type Proxy struct {
	ListenAddr string
	TargetAddr string

	Transformer codec.Transformer
}

func NewProxy(listenAddr, targetAddr string, transformer codec.Transformer) *Proxy {
	return &Proxy{
		ListenAddr:  listenAddr,
		TargetAddr:  targetAddr,
		Transformer: transformer,
	}
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", p.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Proxy listening on %s, forwarding to %s", p.ListenAddr, p.TargetAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go p.handleConnection(clientConn)
	}
}

func (p *Proxy) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	targetConn, err := net.Dial("tcp", p.TargetAddr)
	if err != nil {
		log.Printf("Error connecting to target: %v", err)
		return
	}
	defer targetConn.Close()

	go p.proxy(clientConn, targetConn)
	p.proxy(targetConn, clientConn)
}

func (p *Proxy) proxy(dst, src net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v", err)
			}
			return
		}

		buf = buf[:n]

		buf = p.Transformer(buf)

		_, err = dst.Write(buf)
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			return
		}
	}
}
