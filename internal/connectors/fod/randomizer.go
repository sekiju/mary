package fod

// Randomizer represents a simple random number generator.
type Randomizer struct {
	next int
}

// Constants for the random number generator parameters.
const (
	paramA  = 1103515245
	paramB  = 12345
	randMax = 32767
)

// NewRandomizer creates a new instance of the Randomizer.
func NewRandomizer(seed string) *Randomizer {
	return &Randomizer{next: strToInt(seed)}
}

// Rand generates a random number in the specified range.
func (r *Randomizer) Rand(e int) int {
	if e != 0 {
		e++
		return r._nextInt() / (randMax/e + 1)
	}

	return r._nextInt()
}

func (r *Randomizer) Shuffle(e []int) []int {
	s := len(e)
	o := make([]int, s)
	copy(o, e)

	for t, n, i := 0, 0, 0; n < s; n++ {
		t = r.Rand(s - 1)
		i = o[t]
		o[t] = o[n]
		o[n] = i
	}

	return o
}

// _nextInt generates the next random integer.
func (r *Randomizer) _nextInt() int {
	r.next = (r.next*paramA + paramB) % (randMax + 1)
	return r.next
}

// strToInt converts a string to an integer.
func strToInt(str string) int {
	var o int

	for i := 0; i < len(str); i += 2 {
		o += int(str[i])<<8 | int(str[i+1])
	}

	return o
}
