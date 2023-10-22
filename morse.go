package morse

import (
	"strings"
	"time"
)

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

// InvalidCharacterError is returned when a character is not representable in
// morse code.
type InvalidCharacterError struct {
	Char rune
}

// Error implements the error interface.
func (e InvalidCharacterError) Error() string {
	return "morse: unrepresentable character \"" + string(e.Char) + "\""
}

// Send encodes a text message over the telegraph's Pin. If the message
// contains unsupported characters, an InvalidCharacterError is returned.
func (t *Telegraph) Send(message string) error {
	for _, char := range message {
		_, err := getMorse(char) // Check for unsupported characters.
		if err != nil {
			return err
		}
	}
	for _, char := range message {
		err := t.sendSingle(char)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Telegraph) sendSingle(char rune) error {
	code, err := getMorse(char)
	if err != nil {
		return err
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
	return nil
}

// Code represents a morse code sequence. May be multiple or single characters.
type Code struct {
	c []string
}

// Encode encodes a utf8 text message into the Code type.
func (c *Code) Encode(message string) error {
	for _, char := range message {
		code, err := getMorse(char)
		if err != nil {
			return err
		}
		c.c = append(c.c, code)
	}
	return nil
}

// String returns the morse code sequence as a string. Characters are separated
// by a space, words are separated by 2 spaces.
func (c *Code) String() string {
	return strings.Join(c.c, " ")
}

// LetterCode returns a morse code sequence for the given character.
// If the character is not representable in morse code an InvalidCharacterError is returned.
func LetterCode(char rune) (Code, error) {
	code, err := getMorse(char)
	if err != nil {
		return Code{}, err
	}
	return Code{c: []string{code}}, nil
}

func getMorse(char rune) (string, error) {
	if char > 255 {
		return "", InvalidCharacterError{Char: char} // UTF8 unsupported.
	}
	if char >= 'a' && char <= 'z' {
		// Is lower case, we convert to upper case.
		char -= 32
	}
	code := morseASCIITable[char]
	if code == "" {
		return "", InvalidCharacterError{Char: rune(char)}
	}
	return code, nil
}

var morseASCIITable = [256]string{
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
