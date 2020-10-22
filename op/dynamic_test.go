package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestDynamic_Storex(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := Dynamic{Target: []core.Sequenceable{l}, Emphasis: "++"}
	if got, want := d.Storex(), "dynamic('++',sequence('A B'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestCheckDynamic(t *testing.T) {
	type args struct {
		emphasis string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"piano", args{"-"}, false},
		{"normal", args{"o"}, false},
		{"forte", args{"+++"}, false},
		{"bogus", args{"~"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckDynamic(tt.args.emphasis); (err != nil) != tt.wantErr {
				t.Errorf("CheckDynamic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
