package main

type Channel struct {
	Messages        []string
	VisibleMessages []string
	MessageOffset   int
	Closable        bool
}

func NewChannel(closable bool) *Channel {
	newChannel := Channel{Closable: closable}
	return &newChannel
}

func (c *Channel) SetMessageItems(height, width int) {

	if len(c.Messages) == 0 {
		return
	}

	//calculate how many messages can fit
	var currentMessages []string

	availableHeight := height
	availableWidth := width

	for i := len(c.Messages) - c.MessageOffset - 1; i >= 0; i-- {
		availableHeight -= len(c.Messages[i]) / availableWidth
		if len(c.Messages[i])%availableWidth > 0 {
			availableHeight--
		}

		if availableHeight < 1 {
			break
		}

		currentMessages = append(currentMessages, c.Messages[i])
	}

	//reverse message order
	var reversedMessages []string
	for i := len(currentMessages) - 1; i >= 0; i-- {
		reversedMessages = append(reversedMessages, currentMessages[i])
	}

	c.VisibleMessages = reversedMessages
}

func (c *Channel) ListMove(y int) {
	if c.MessageOffset+y < 0 || c.MessageOffset+y > len(c.Messages)-1 {
		return
	}
	c.MessageOffset += y
	renderScreen()
}
