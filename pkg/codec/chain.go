package codec

import "errors"

type Chain struct {
	Codecs []Codec
}

func NewChain(codecs []Codec) (*Chain, error) {
	if len(codecs) == 0 {
		return nil, errors.New("at least one codec is required")
	}
	return &Chain{Codecs: codecs}, nil
}

func (c *Chain) NewEncoder() Transformer {
	encoders := make([]Transformer, len(c.Codecs))
	for i, codec := range c.Codecs {
		encoders[i] = codec.NewEncoder()
	}

	return func(data []byte) []byte {
		result := data
		// Apply encoders in forward order
		for _, encoder := range encoders {
			result = encoder(result)
		}
		return result
	}
}

func (c *Chain) NewDecoder() Transformer {
	decoders := make([]Transformer, len(c.Codecs))
	for i, codec := range c.Codecs {
		decoders[i] = codec.NewDecoder()
	}

	return func(data []byte) []byte {
		result := data
		// Apply decoders in reverse order
		for i := len(decoders) - 1; i >= 0; i-- {
			result = decoders[i](result)
		}
		return result
	}
}
