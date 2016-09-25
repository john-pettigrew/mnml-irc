package main

type InputBuffer struct {
	Pos      int
	Contents string
}

func (i *InputBuffer) Clear() {
	i.Pos = 0
	i.Contents = ""
}

func (i *InputBuffer) CursorMove(value int) {

	//Edge cases
	if i.Pos+value > len(i.Contents)-1 || i.Pos+value < 0 {
		return
	}

	i.Pos += value
}

func (i *InputBuffer) Type(input rune) {

	//Handle directions
	if input == rune(19) {
		if i.Pos == 0 {
			return
		}
		i.Pos -= 1
		return
	}
	if input == rune(18) {
		if i.Pos == len(i.Contents)-1 {
			return
		}
		i.Pos += 1
		return
	}

	//Handle backspace
	if input == rune(127) {
		if i.Pos == 0 {
			return
		}
		i.Contents = i.Contents[:i.Pos-1] + i.Contents[i.Pos:]
		i.Pos -= 1
		return
	}

	//Handle regular input
	i.Contents = i.Contents[:i.Pos] + string(input) + i.Contents[i.Pos:]
	i.Pos += 1
}
