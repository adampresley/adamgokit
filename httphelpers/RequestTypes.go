package httphelpers

type RequestTypes interface {
	int | []int |
		int32 | []int32 |
		int64 | []int64 |
		uint | []uint |
		uint32 | []uint32 |
		uint64 | []uint64 |
		float32 | []float32 |
		float64 | []float64 |
		string | []string |
		bool
}
