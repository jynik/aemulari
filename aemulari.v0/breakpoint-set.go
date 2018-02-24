package aemulari

import "sort"

// A set of breakpoints, accessible by unique ID or the address at which
// they're placed. In the future, I'd like to have conditional breakpoints,
// hence multiple bps at one address.
type breakpointSet struct {
	nextId int                      // Monotonically increasing
	byAddr map[uint64][]*Breakpoint // Address -> BP's at that address
	byID   map[int]*Breakpoint      // ID -> BP
}

func (bps *breakpointSet) initialize() {
	bps.nextId = 1
	bps.byAddr = make(map[uint64][]*Breakpoint)
	bps.byID = make(map[int]*Breakpoint)
}

func (bps *breakpointSet) add(addr uint64) Breakpoint {
	id := bps.nextId
	bps.nextId++

	bp := newBreakpoint(id, addr)

	if list, present := bps.byAddr[addr]; !present {
		bps.byAddr[addr] = []*Breakpoint{&bp}
	} else {
		bps.byAddr[addr] = append(list, &bp)
	}

	bps.byID[id] = &bp

	return bp
}

// Returns true if any BPs at `addr` were triggered
func (bps *breakpointSet) process(addr uint64) bool {
	var anyTriggered bool = false

	for _, bp := range bps.byID {
		triggered := bp.hit(addr)
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
func (bps *breakpointSet) removeAll() {
	bps.initialize()
}

// Remove all breakpoints at the specified address
func (bps *breakpointSet) removeAllAt(addr uint64) {
	if byAddr, present := bps.byAddr[addr]; present {
		for _, bp := range byAddr {
			delete(bps.byID, bp.ID)
		}

		delete(bps.byAddr, addr)
	}
}

// Remove the breakpoint associated with the specified Id
func (bps *breakpointSet) remove(id int) {
	if bp, present := bps.byID[id]; present {
		if addrList, present := bps.byAddr[bp.Address]; present {
			if len(addrList) == 1 && addrList[0].ID == id {
				// Just zap the list, it's the last breakpoint at this address
				delete(bps.byAddr, bp.Address)
			} else {
				// Find and delete the breakpoint
				for i, bp := range addrList {
					if bp.ID == id {
						// https://github.com/golang/go/wiki/SliceTricks
						bps.byAddr[bp.Address] = append(addrList[:i], addrList[i+1:]...)
						break
					}
				}
			}
		}

		delete(bps.byID, id)
	}
}

// Get all breakpoints, sorted by ID
func (bps breakpointSet) get() BreakpointList {
	max := len(bps.byID)

	var ids []int = make([]int, max, max)
	var ret BreakpointList = make(BreakpointList, max, max)

	i := 0
	for id := range bps.byID {
		ids[i] = id
		i++
	}

	sort.Ints(ids)
	for i, id := range ids {
		ret[i] = *(bps.byID[id])
	}

	return ret
}

// Return list of breakpoints at the specified address, sorted by ID
func (bps breakpointSet) getAllAt(addr uint64) BreakpointList {
	ptrs, present := bps.byAddr[addr]
	if !present {
		return BreakpointList{}
	}

	max := len(ptrs)
	ids := make([]int, max, max)
	ret := make(BreakpointList, max, max)

	i := 0
	for _, bp := range ptrs {
		ids[i] = bp.ID
		i++
	}

	sort.Ints(ids)

	i = 0
	for _, id := range ids {
		ret[i] = *(bps.byID[id])
		i++
	}

	return ret
}
