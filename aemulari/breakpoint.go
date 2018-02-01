package aemulari

import (
	"fmt"
	"sort"
)

type Breakpoint struct {
	Id      int
	Count   uint
	Address uint64
	state   breakpointState
}

type BreakpointList []Breakpoint

// Maps an address to one or more breakpoints at that address.
// In the future, I'd like to have conditional breakpoints, hence multiple bps at one address.
type Breakpoints struct {
	nextId int                      // Monotonically increasing
	byAddr map[uint64][]*Breakpoint // Address -> BP's at that address
	byId   map[int]*Breakpoint      // ID -> BP
}

type breakpointState int

const (
	breakpointInvalid = iota
	breakpointInactive
	breakpointArmed
	breakpointTriggered
	breakpointMax
)

func newBreakpoint(id int, addr uint64) Breakpoint {
	var b Breakpoint

	b.Address = addr
	b.Id = id
	b.Reset()

	return b
}

func (b *Breakpoint) Reset() {
	b.Count = 0
	b.Enable()
}

func (b *Breakpoint) Disable() {
	b.state = breakpointInactive
}

func (b *Breakpoint) Enable() {
	b.state = breakpointArmed
}

func (b *Breakpoint) Enabled() bool {
	return b.state != breakpointInactive
}

//
func (b *Breakpoint) Hit(addr uint64) bool {
	if addr != b.Address {
		return false
	}

	if b.state <= breakpointInvalid || b.state >= breakpointMax {
		log.Warning(fmt.Sprintf("BP is in valid state (%d)", b.state))
		return false
	}

	b.Count++

	if b.state == breakpointArmed {
		b.state = breakpointTriggered
		return true
	}

	return false
}

func (b Breakpoint) String() string {
	// FIXME address format string should be arch-dependent
	return fmt.Sprintf("Breakpoint %2d: 0x%08x, Hit count = %d", b.Id, b.Address, b.Count)
}

// Returns true if any of the breakpoints in the provided list are enabled
func (bpl BreakpointList) Enabled() bool {
	for _, b := range bpl {
		if b.Enabled() {
			return true
		}
	}
	return false
}

func (bps *Breakpoints) Initialize() {
	bps.nextId = 1
	bps.byAddr = make(map[uint64][]*Breakpoint)
	bps.byId = make(map[int]*Breakpoint)
}

func (bps *Breakpoints) Add(addr uint64) Breakpoint {
	id := bps.nextId
	bps.nextId++

	bp := newBreakpoint(id, addr)

	if list, present := bps.byAddr[addr]; !present {
		bps.byAddr[addr] = []*Breakpoint{&bp}
	} else {
		bps.byAddr[addr] = append(list, &bp)
	}

	bps.byId[id] = &bp

	return bp
}

// Returns true if any BPs at `addr` were triggered
func (bps *Breakpoints) Process(addr uint64) bool {
	var anyTriggered bool = false

	for _, bp := range bps.byId {
		triggered := bp.Hit(addr)
		anyTriggered = anyTriggered || triggered

		if bp.Address != addr {
			if bp.state == breakpointTriggered {
				/* Re-arm breakpoint. This must be done after we leave the address,
				* else we'll get "stuck" */
				bp.Enable()
			}
		}
	}

	return anyTriggered
}

// Remove all breakpoints
func (bps *Breakpoints) RemoveAll() {
	bps.Initialize()
}

// Remove all breakpoints at the specified address
func (bps *Breakpoints) RemoveAllAt(addr uint64) {
	if byAddr, present := bps.byAddr[addr]; present {
		for _, bp := range byAddr {
			delete(bps.byId, bp.Id)
		}

		delete(bps.byAddr, addr)
	}
}

// Remove the breakpoint associated with the specified Id
func (bps *Breakpoints) Remove(id int) {
	if bp, present := bps.byId[id]; present {
		if addrList, present := bps.byAddr[bp.Address]; present {
			if len(addrList) == 1 && addrList[0].Id == id {
				// Just zap the list, it's the last breakpoint at this address
				delete(bps.byAddr, bp.Address)
			} else {
				// Find and delete the breakpoint
				for i, bp := range addrList {
					if bp.Id == id {
						// https://github.com/golang/go/wiki/SliceTricks
						bps.byAddr[bp.Address] = append(addrList[:i], addrList[i+1:]...)
						break
					}
				}
			}
		}

		delete(bps.byId, id)
	}
}

// Get all breakpoints, sorted by ID
func (bps Breakpoints) Get() BreakpointList {
	max := len(bps.byId)

	var ids []int = make([]int, max, max)
	var ret BreakpointList = make(BreakpointList, max, max)

	i := 0
	for id := range bps.byId {
		ids[i] = id
		i++
	}

	sort.Ints(ids)
	for i, id := range ids {
		ret[i] = *(bps.byId[id])
	}

	return ret
}

// Return list of breakpoints at the specified address, sorted by ID
func (bps Breakpoints) GetAt(addr uint64) BreakpointList {
	ptrs, present := bps.byAddr[addr]
	if !present {
		return BreakpointList{}
	}

	max := len(ptrs)
	ids := make([]int, max, max)
	ret := make(BreakpointList, max, max)

	i := 0
	for _, bp := range ptrs {
		ids[i] = bp.Id
		i++
	}

	sort.Ints(ids)

	// FIXME Yea, this is rather wasteful... :shrug:
	i = 0
	for _, id := range ids {
		ret[i] = *(bps.byId[id])
		i++
	}

	return ret
}
