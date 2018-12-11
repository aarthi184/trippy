package spinner

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func init() {
	var seed int
	host, err := os.Hostname()
	if err == nil {
		for _, n := range host {
			seed = seed + int(n)
		}
	} else {
		fmt.Println("Error in getting host name:", err)
	}
	// Randomizing the seed using machine hostname & current time
	// So that all machines don't have the same seed if started at the same time
	rand.Seed(time.Now().UnixNano() + int64(seed))
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
