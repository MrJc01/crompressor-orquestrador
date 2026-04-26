package hamming

import "testing"

func TestDistance(t *testing.T) {
	tests := []struct {
		name     string
		hashA    uint64
		hashB    uint64
		expected int
	}{
		{"Iguais", 0b101010, 0b101010, 0},
		{"Diferenca de 1 bit", 0b101010, 0b101011, 1},
		{"Diferenca total", 0xFFFFFFFFFFFFFFFF, 0x0000000000000000, 64},
		{"Diferenca parcial", 0b11110000, 0b00001111, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Distance(tt.hashA, tt.hashB); got != tt.expected {
				t.Errorf("Distance() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsSimilar(t *testing.T) {
	// A=1010 (10), B=1011 (11) -> diff 1 bit
	if !IsSimilar(10, 11, 1) {
		t.Errorf("Esperado ser similar com threshold 1")
	}
	if IsSimilar(10, 11, 0) {
		t.Errorf("Nao esperado ser similar com threshold 0")
	}
}
