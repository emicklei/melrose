package dsl

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{
			"valid",
			`sequence('C D E F G A B C5')`,
			false,
		},
		{
			"unknown function",
			`sequenc('C')`,
			true,
		},
		{
			"unknown note",
			`sequence('X')`,
			true,
		},
		{
			"unterminated string",
			`sequence('A)`,
			true,
		},
		{
			"no function",
			`('A')`,
			false,
		},
		{
			"wrong arg type",
			`octave('A','B')`,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if err != nil {
					t.Logf("Received expected error: %v", err)
				}
			}
		})
	}
}
