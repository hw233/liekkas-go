package httputil

type HeaderField struct {
	K, V string
}

func NewHeaderField(K, V string) *HeaderField {
	return &HeaderField{
		K: K,
		V: V,
	}
}
