package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"github.com/uia-worker/is105sem03/mycrypt"
	"github.com/liahra/minyr/yr"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)

					msg := string(dekryptertMelding)

					switch {
  				        case msg == "ping":
						krypterPong := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(string(krypterPong)))
					case strings.Contains(msg, "Kjevik"):
						output, err := yr.CelsiusToFahrenheitLine(msg)
						if err != nil {
							log.Println(err)
						}

						log.Println(output)
						kryptertOutput := mycrypt.Krypter([]rune(output), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(string(kryptertOutput)))
					default:
						_, err = c.Write(buf[:n])
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
