package codec

import "testing"

func TestEncoder(t *testing.T) {
	c, err := NewEncoder(1, 2, []byte("3")).Encode()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	t.Logf("content: %v", c)
}
