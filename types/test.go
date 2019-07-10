package types

// Test is the struct for testing
type Test struct {
	St  string  `json:"st"`
	B   bool    `json:"b"`
	C   uint8   `json:"c"`
	S   int16   `json:"s"`
	Us  uint16  `json:"us"`
	L   int32   `json:"l"`
	Ul  uint32  `json:"ul"`
	Ll  int64   `json:"ll"`
	Ull uint64  `json:"ull"`
	F   float32 `json:"f"`
	D   float64 `json:"d"`
}
