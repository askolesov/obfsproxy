package codec

// None is a codec that does not modify the data
type None struct{}

// NewNone creates a new None codec
func NewNone() *None {
	return &None{}
}

// NewEncoder returns a transformer that returns data unchanged
func (n *None) NewEncoder() Transformer {
	return func(data []byte) []byte {
		return data
	}
}

// NewDecoder returns a transformer that returns data unchanged
func (n *None) NewDecoder() Transformer {
	return func(data []byte) []byte {
		return data
	}
}
