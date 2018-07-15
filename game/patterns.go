package game

// NumKeys specifies the keymode
const NumKeys = 4

// Note translates the noteskin index atlas
type Note struct {
	Rune rune
}

const (
	// Receptor represented by a character
	Receptor rune = 'R'
	// Red represented by a character
	Red = 'r'
	// Blue represented by a character
	Blue = 'b'
	// Yellow represented by a character
	Yellow = 'y'
)

// MakeStream generates a stream of non-repeating notes
func (c *Character) MakeStream(len int) []int {
	last := -1
	notes := make([]int, len)
	for i := range notes {
		notes[i] = c.PatternRNG.Intn(NumKeys)
		for notes[i] == last {
			notes[i] = c.PatternRNG.Intn(NumKeys)
		}
		last = notes[i]
	}
	return notes
}
