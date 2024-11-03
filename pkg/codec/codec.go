package codec

type Codec interface {
	NewEncoder() Transformer
	NewDecoder() Transformer
}

type Transformer func(data []byte) []byte
