package pkg

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/askolesov/obfsproxy/pkg/codec"
)

type Proxy struct {
	ListenAddr string
	TargetAddr string

	IsServer bool
	Codec    codec.Codec
}

func NewProxy(listenAddr, targetAddr string, isServer bool, codec codec.Codec) *Proxy {
	return &Proxy{
		ListenAddr: listenAddr,
		TargetAddr: targetAddr,
		IsServer:   isServer,
		Codec:      codec,
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
		log.Printf("Accepted connection from %s", clientConn.RemoteAddr())
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
	log.Printf("Connected to target %s", targetConn.RemoteAddr())
	defer targetConn.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	if p.IsServer {
		go func() {
			defer wg.Done()
			p.proxy(clientConn, targetConn, p.Codec.NewDecoder())
		}()
		go func() {
			defer wg.Done()
			p.proxy(targetConn, clientConn, p.Codec.NewEncoder())
		}()
	} else {
		go func() {
			defer wg.Done()
			p.proxy(clientConn, targetConn, p.Codec.NewEncoder())
		}()
		go func() {
			defer wg.Done()
			p.proxy(targetConn, clientConn, p.Codec.NewDecoder())
		}()
	}

	wg.Wait()
}

func (p *Proxy) proxy(src, dst net.Conn, t codec.Transformer) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			data := buf[:n]
			data = t(data)

			_, writeErr := dst.Write(data)
			if writeErr != nil {
				log.Printf("Error writing to connection: %v", writeErr)
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v", err)
			}
			return
		}
	}
}
