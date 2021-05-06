package aemulari

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
//	"time"
)

// Interoperates with gdbghidra plugin to synchronize
// aemulari debug state with Ghidra disassembly and decompiler cursor
// locations.
//
// gdbghidra: https://github.com/Comsecuris/gdbghidra
type ghidraBridge struct {

	ghidraConn net.Conn	// Commands we send to the Ghidra plugin
	isConnected bool	// Is ghidraConn an active connection?

	eventListener *net.TCPListener	// Async events from Ghidra to us
	isListening bool				// Is eventConn actively listening?
	shutdownListener chan bool		/* Any value sent o this channel
									 * will shut down the listerner and
									 * any open connection. */
}

// Connect to ghidraBridge and start up listening service for
// events originating from Ghidra. The default ports for the plugin
// are used, and 127.0.0.1 is assumed.
func (gb *ghidraBridge) Connect() error {
	//var localAddr net.TCPAddr
	var err error

	if gb.isConnected {
		gb.Disconnect()
	}

	gb.ghidraConn, err = net.Dial("tcp", "127.0.0.1:2305")
	if err != nil {
		return err
	}
	gb.isConnected = true

	//localAddr.IP = net.IPv4(127, 0, 0, 1)
	//localAddr.Port = 2306

	//gb.eventListener, err = net.ListenTCP("tcp", &localAddr)
	//if err != nil {
	//	gb.ghidraConn.Close()
	//	return err
	//}
	//gb.isListening = true
	//go gb.listenerLoop()

	return nil
}

func (gb *ghidraBridge) SetCursorAddress(addr uint64) error {
	if !gb.isConnected {
		return nil
	}

	if msg, err := updateCursorMessage(addr); err != nil {
		return err
	} else {
		msg = append(msg, '\n')
		_, err = gb.ghidraConn.Write(msg)
		return err
	}

	return nil
}

func (gb *ghidraBridge) Disconnect() {
	if gb.isListening {
		gb.shutdownListener <- true
		close(gb.shutdownListener)
	}

	if gb.isConnected {
		gb.isConnected = false
		gb.ghidraConn.Close()
	}
}

//func (gb *ghidraBridge) listenerLoop() {
//	shutdown := false
//
//	fmt.Println("Starting listernaerloop")
//
//	for !shutdown {
//		fmt.Println("Checking channel")
//		_, shutdown := <-gb.shutdownListener
//		if shutdown {
//			continue
//		}
//		fmt.Println("Checked channel, now Accept()'ing")
//
//		gb.eventListener.SetDeadline(time.Now().Add(1 * time.Second))
//		conn, err := gb.eventListener.Accept()
//		if (err != nil) {
//			fmt.Println("%s\n", err)
//			time.Sleep(250 * time.Millisecond)
//			continue
//		}
//
//		fmt.Printf("Got event connection from %s\n", conn.RemoteAddr().String())
//
//		shutdown = gb.eventHandler(conn)
//	}
//
//	gb.eventListener.Close()
//	gb.isListening = false
//}
//
//func (gb *ghidraBridge) eventHandler(conn net.Conn) bool {
//	return false
//}


type bridgeMessage struct {
	Type string `json:"type"`
	Data []map[string]string `json:"data"`	// Unclear why the authors made this an array of maps...
}

func encodeBridgeMessage(msgType string, fields map[string]string) ([]byte, error) {
	var msg bridgeMessage
	msg.Type = msgType
	msg.Data = append(msg.Data, fields)
	return json.Marshal(msg)
}

func helloMessage(localAddr string) ([]byte, error) {

	sep := strings.LastIndex(localAddr, ":")
	if sep < 0 {
		msg := "Argument not in expected <address>:<port> format: %s"
		err := fmt.Errorf(msg, localAddr)
		return []byte{}, err
	}

	ipStr := localAddr[0:sep]
	portStr := localAddr[sep+1:]

	fields := make(map[string]string)

	// TODO: Adjust later if used by the plugin. They don't seem to be
	// Unusued by plugin, just
	fields["arch"] = "ARM"
	fields["endian"] = "little"
	fields["answer_ip"] = ipStr
	fields["answer_port"] = portStr

	return encodeBridgeMessage("HELLO", fields)
}

func updateCursorMessage(address uint64) ([]byte, error) {
	fields := make(map[string]string)
	fields["address"] = fmt.Sprintf("%d", address)

	// Don't worry about this for now.
	// TODO: Determine if useful to adjust based upon mappings in emulator.
	fields["relocate"] = "0"

	return encodeBridgeMessage("CURSOR", fields)
}

//func main() {
//	var gb ghidraBridge
//
//	err := gb.Connect()
//	if err != nil {
//		panic(err)
//	}
//
//	err = gb.SetCursorAddress(0x1924)
//	if err != nil {
//		panic(err)
//	}
//	gb.Disconnect()
//}
