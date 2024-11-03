package codec

import "errors"

var _ Codec = &Xorer{}

type Xorer struct {
	Key []byte
}

func NewXorer(key []byte) (*Xorer, error) {
	if len(key) == 0 {
		return nil, errors.New("key must not be empty")
	}

	return &Xorer{Key: key}, nil
}

func (x *Xorer) NewEncoder() Transformer {
	keyIndex := 0

	return func(data []byte) []byte {
		result := make([]byte, len(data))
		for i, b := range data {
			result[i] = b ^ x.Key[keyIndex]
			keyIndex = (keyIndex + 1) % len(x.Key)
		}
		return result
	}
}

func (x *Xorer) NewDecoder() Transformer {
	keyIndex := 0

	return func(data []byte) []byte {
		result := make([]byte, len(data))
		for i, b := range data {
			result[i] = b ^ x.Key[keyIndex]
			keyIndex = (keyIndex + 1) % len(x.Key)
		}
		return result
	}
}
