package hash

import (
	"testing"
)

func Test_Sha512(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"Test value 'I like apples' returns expected. ", "I like apples", "921fde90136f6f6a0edf3cc4fced9533eff1133f945f9cbe2939d9c4d1135b65d5b2ce9bdf0d32934a0bd8e6302c9141a90a3888a1cc4887a553011008a5e994"},
		{"Test empty value returns expected.", "", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"Test value with special characters returns expected. ", ":asdf#", "5ab74ab02ccf9d7617534419344415e93bc56622268ae781df98c344662bcc483e7818d31ed341964106942aa59a7fdd4aa267264229fb767682ec1ff4510729"},
		{"Test unicode character ă returns expected", "ă", "cd15b4ad6d85be0404acc0f921677ee6206893cee71642575660e8186a20d57b734c22e8aa57c331b7bf9e1fca5bb4f2b7f37d13f04ce9297717ee05fd82c6fc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sha512(tt.value)
			if got != tt.want {
				t.Errorf("Sha512() got = %v, want %v", got, tt.want)
			}
		})
	}
}
