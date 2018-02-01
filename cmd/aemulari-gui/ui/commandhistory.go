package ui

import "fmt"

type CommandHistoryItem struct {
	status string
	line   string
}

type CommandHistory struct {
	index   int
	entries []CommandHistoryItem
}

func (i CommandHistoryItem) String() string {
	return fmt.Sprintf("%s %s\n", i.status, i.line)
}

func (c *CommandHistory) Clear() {
	c.index = 0
	c.entries = []CommandHistoryItem{}
}

func (c CommandHistory) Size() int {
	return len(c.entries)
}

func (c *CommandHistory) PrevEntry() CommandHistoryItem {
	var ret CommandHistoryItem

	if c.index > 0 {
		c.index--
	}

	if c.Size() > 0 {
		return c.entries[c.index]
	}

	return ret

}

func (c *CommandHistory) PrevEntryLine() string {
	return c.PrevEntry().line
}

func (c *CommandHistory) NextEntry() CommandHistoryItem {
	var ret CommandHistoryItem

	size := c.Size()
	if size <= 0 {
		return ret
	}

	if c.index < (size - 1) {
		c.index++
	} else if c.index >= size {
		c.index = size - 1
	}

	ret = c.entries[c.index]

	return ret
}

func (c *CommandHistory) NextEntryLine() string {
	return c.NextEntry().line
}

// Append command to history if it isn't a duplicate of the previous command
func (c *CommandHistory) Append(status, command string) {
	size := c.Size()
	if size == 0 || c.entries[size-1].line != command {
		c.entries = append(c.entries, CommandHistoryItem{status, command})
	}
	c.ResetIndex()
}

func (c CommandHistory) EntryString(i int) string {
	if i >= 0 && i < c.Size() {
		return c.entries[i].String()
	} else {
		return ""
	}
}

func (c *CommandHistory) ResetIndex() {
	c.index = c.Size()
}
