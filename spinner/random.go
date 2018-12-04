package spinner

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	// TODO: Figure out if seed needs to be changed once in a while
	rand.Seed(time.Now().UnixNano())
}

var (
	errMinNotLess = errors.New("min is not less than max")
)

// randInt returns a random number in the interval [min,max]
// both min and max are included
func randInt(min, max int) (int, error) {
	// rand.Intn(n) panics if n is 0
	// return err if max is less than or equal to min
	if max <= min {
		return -1, errMinNotLess
	}

	// rand.Intn(n) returns a random number in the interval [0,n)
	// adding 1 to include n
	return min + rand.Intn((max-min)+1), nil
}
