package op

import (
	"testing"

	"github.com/emicklei/melrose/core"
)

func TestDynamic_Storex(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := Dynamic{Target: []core.Sequenceable{l}, Emphasis: core.On("++")}
	if got, want := d.Storex(), "dynamic('++',sequence('A B'))"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	{
		l := core.MustParseSequence("A B")
		d := Dynamic{Target: []core.Sequenceable{l}, Emphasis: core.On(54)}
		if got, want := d.Storex(), "dynamic(54,sequence('A B'))"; got != want {
			t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
		}
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
func TestDynamic_S(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := Dynamic{Target: []core.Sequenceable{l}, Emphasis: core.On("++")}
	if got, want := d.S().Storex(), "sequence('A++ B++')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	d = Dynamic{Target: []core.Sequenceable{l}, Emphasis: core.On(127)}
	if got, want := d.S().Storex(), "sequence('A+++++ B+++++')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	d = Dynamic{Target: []core.Sequenceable{l}, Emphasis: core.On("invalid")}
	if got, want := d.S().Storex(), "sequence('')"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func TestDynamic_Replaced(t *testing.T) {
	l := core.MustParseSequence("A B")
	d := Dynamic{Target: []core.Sequenceable{l}, Emphasis: core.On("++")}
	if core.IsIdenticalTo(d, l) {
		t.Error("should not be identical")
	}
	if !core.IsIdenticalTo(d.Replaced(l, core.EmptySequence).(Dynamic).Target[0], core.EmptySequence) {
		t.Error("not replaced")
	}
	if !core.IsIdenticalTo(d.Replaced(d, l), l) {
		t.Error("should be replaced by l")
	}
}

func TestDynamic_ToNote(t *testing.T) {
	n := core.MustParseNote("C")
	d := Dynamic{Target: []core.Sequenceable{n}, Emphasis: core.On("++")}
	not, err := d.ToNote()
	if err != nil {
		t.Fatal(err)
	}
	if not.Velocity != 80 {
		t.Errorf("got %d want %d", not.Velocity, 80)
	}

	d = Dynamic{Target: []core.Sequenceable{n}, Emphasis: core.On(42)}
	not, err = d.ToNote()
	if err != nil {
		t.Fatal(err)
	}
	if not.Velocity != 42 {
		t.Errorf("got %d want %d", not.Velocity, 42)
	}

	d = Dynamic{Target: []core.Sequenceable{n}, Emphasis: core.On("invalid")}
	_, err = d.ToNote()
	if err == nil {
		t.Fatal("error expected")
	}

	d = Dynamic{Target: []core.Sequenceable{}, Emphasis: core.On(42)}
	_, err = d.ToNote()
	if err == nil {
		t.Fatal("error expected")
	}

	d = Dynamic{Target: []core.Sequenceable{core.MustParseSequence("C")}, Emphasis: core.On(42)}
	_, err = d.ToNote()
	if err == nil {
		t.Fatal("error expected")
	}

	d = Dynamic{Target: []core.Sequenceable{failingNoteConvertable{}}, Emphasis: core.On(42)}
	_, err = d.ToNote()
	if err == nil {
		t.Fatal("error expected")
	}
}
