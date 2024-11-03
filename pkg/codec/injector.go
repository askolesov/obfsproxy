package codec

import (
	"errors"

	"golang.org/x/exp/rand"
)

// Codec

var _ Codec = &Injector{}

type Injector struct {
	Seed       uint64
	Redundancy int

	rate int
}

func NewInjector(seed uint64, redundancy int) (*Injector, error) {
	if redundancy < 0 || redundancy > 1000 {
		return nil, errors.New("redundancy must be between 0 and 1000")
	}

	rate := rate(redundancy)

	return &Injector{Seed: seed, Redundancy: redundancy, rate: rate}, nil
}

func rate(redundancy int) int {
	return 100 * redundancy / (100 + redundancy)
}

func (inj *Injector) NewEncoder() Transformer {
	rng := rand.New(rand.NewSource(inj.Seed))

	return func(input []byte) []byte {
		n := len(input)
		output := make([]byte, 0, n)
		for i := 0; i < n; i++ {
			if rng.Intn(100) < inj.rate { // inject
				output = append(output, byte(rng.Intn(256)))
				i-- // roll back
			} else {
				output = append(output, input[i])
			}
		}
		return output
	}
}

func (inj *Injector) NewDecoder() Transformer {
	rng := rand.New(rand.NewSource(inj.Seed))

	return func(input []byte) []byte {
		n := len(input)
		output := make([]byte, 0, n)
		for i := 0; i < n; i++ {
			if rng.Intn(100) < inj.rate { // discard the random byte
				_ = rng.Intn(256)
			} else {
				output = append(output, input[i])
			}
		}
		return output
	}
}
