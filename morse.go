package morse

import "time"

// Pin is an interface that allows the telegraph to be used with different
// hardware. A Pin can be implemented with a GPIO pin or even a buzzer.
type Pin func(logicLevel bool)

// Telegraph implements morse encoding over a Pin type.
type Telegraph struct {
	dot time.Duration
	pin Pin
}

// NewTelegraph returns a new telegraph instance ready for use.
func NewTelegraph(dot time.Duration, pin Pin) *Telegraph {
	if dot <= 0 {
		panic("dot duration must be greater than 0")
	}
	return &Telegraph{dot: dot, pin: pin}
}

// Send encodes a message over the telegraph's Pin.
func (t *Telegraph) Send(message string) {
	for _, char := range message {
		if char > 255 {
			continue // Ignore non-ASCII characters.
		}
		t.sendASCII(byte(char))
	}
}

func (t *Telegraph) sendASCII(char byte) {
	if char >= 'a' && char <= 'z' {
		// Is lower case, we convert to upper case.
		char -= 32
	}
	code := morseTable[char]
	if code == "" {
		return // Ignore unknown characters.
	}
	for _, symbol := range code {
		if symbol != ' ' {
			t.pin(true)
		}
		switch symbol {
		case '*':
			time.Sleep(t.dot)
		case '-': // Dashes or "dahs" are 3 times longer than dots or "dits".
			time.Sleep(t.dot * 3)
		case ' ': // Spaces between words are equal to 7 dots.
			time.Sleep(t.dot * 6)
		}
		t.pin(false)
		time.Sleep(t.dot) // Inter-character gap is 1 dot.
	}
}

var morseTable = [256]string{
	' ': " ",
	'A': "*-",
	'B': "-***",
	'C': "-*-*",
	'D': "-**",
	'E': "*",
	'F': "**-*",
	'G': "--*",
	'H': "****",
	'I': "**",
	'J': "*---",
	'K': "-*-",
	'L': "*-**",
	'M': "--",
	'N': "-*",
	'O': "---",
	'P': "*--*",
	'Q': "--*-",
	'R': "*-*",
	'S': "***",
	'T': "-",
	'U': "**-",
	'V': "***-",
	'W': "*--",
	'X': "-**-",
	'Y': "-*--",
	'Z': "--**",
	'0': "-----",
	'1': "*----",
	'2': "**---",
	'3': "***--",
	'4': "****-",
	'5': "*****",
	'6': "-****",
	'7': "--***",
	'8': "---**",
	'9': "----*",
	'.': "*-*-*-",
	',': "--**--",
	'?': "**--**",
	'"': "*-**-*",
	'/': "-**-*",
}
