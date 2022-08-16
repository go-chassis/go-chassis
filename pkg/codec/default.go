package codec

import (
	"encoding/json"
	"github.com/go-chassis/cari/codec"
)

// StdJson implement standard json codec
type StdJson struct {
}

func newDefault(opts Options) (codec.Codec, error) {
	return &StdJson{}, nil
}
func (s *StdJson) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (s *StdJson) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
