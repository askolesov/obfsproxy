package codec

var _ Codec = &Inverter{}

type Inverter struct{}

func NewInverter() *Inverter {
	return &Inverter{}
}

func (b *Inverter) invert(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = ^b
	}
	return result
}

func (b *Inverter) NewEncoder() Transformer {
	return b.invert
}

func (b *Inverter) NewDecoder() Transformer {
	return b.invert
}
