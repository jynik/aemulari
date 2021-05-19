// Functionality for synchronizing external tools to debugger state
//
// Currently, this is just Ghidra + AddressSync, because that's
// good enough for me! (For now...)
//
// https://github.com/jynik/AddressSync
//

package aemulari

import (
	"fmt"
	"net"
)

type ToolSync struct {
	isInitialized bool

	addrConn net.Conn
	addrBuf []byte
}

// Open connection(s) to external tools
func (ts *ToolSync) Open() error {
	var err error

	ts.addrBuf = make([]byte, 8)

	ts.addrConn, err = net.Dial("udp","127.0.0.1:1080")
	if err != nil {
		return err
	}

	ts.isInitialized = true
	return nil
}

func (ts *ToolSync) Close() {
	if !ts.isInitialized {
		return
	}

	ts.isInitialized = false
	ts.addrConn.Close()
}

// Send the program counter address to any sync'd tools
// If Open() has not been called, this is a No-op that returns no error.
func (ts *ToolSync) SendCurrAddress(addr uint64) error {
	if !ts.isInitialized {
		// Fail gracefully so we can always call this safely sans state test
		return nil
	}

	ts.addrBuf[0] = byte( addr        & 0xff)
	ts.addrBuf[1] = byte((addr >> 8)  & 0xff)
	ts.addrBuf[2] = byte((addr >> 16) & 0xff)
	ts.addrBuf[3] = byte((addr >> 24) & 0xff)
	ts.addrBuf[4] = byte((addr >> 32) & 0xff)
	ts.addrBuf[5] = byte((addr >> 40) & 0xff)
	ts.addrBuf[6] = byte((addr >> 48) & 0xff)
	ts.addrBuf[7] = byte((addr >> 56) & 0xff)

	n, err := ts.addrConn.Write(ts.addrBuf)
	if err != nil {
		return err
	}

	if n != 8 {
		msg := "Tried to send 8 address bytes to external tool, actually sent %d"
		return fmt.Errorf(msg, n)
	}

	return nil;
}
