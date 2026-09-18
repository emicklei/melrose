package core

import (
	"testing"
)

func TestInterval_Value(t *testing.T) {
	i := NewInterval(On(0), On(1), On(1), RepeatFromTo)
	if got, want := i.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	i.Next()
	if got, want := i.Value(), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	i.Next()
	if got, want := i.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInterval_Value_Backwards(t *testing.T) {
	i := NewInterval(On(0), On(1), On(-1), RepeatFromTo)
	if got, want := i.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	i.Next()
	if got, want := i.Value(), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInterval_Strategies(t *testing.T) {
	t.Run("repeat", func(t *testing.T) {
		i := NewInterval(On(0), On(2), On(1), RepeatFromTo)
		for _, want := range []int{0, 1, 2, 0} {
			if got := i.Value(); got != want {
				t.Fatalf("before step got %v want %v", got, want)
			}
			i.Next()
		}
	})

	t.Run("once", func(t *testing.T) {
		i := NewInterval(On(0), On(2), On(1), OnceFromTo)
		for _, want := range []int{0, 1, 2, 2, 2} {
			if got := i.Value(); got != want {
				t.Fatalf("before step got %v want %v", got, want)
			}
			i.Next()
		}
	})

	t.Run("once-two-way", func(t *testing.T) {
		i := NewInterval(On(0), On(2), On(1), OnceFromToFrom)
		for _, want := range []int{0, 1, 2, 1, 0} {
			if got := i.Value(); got != want {
				t.Fatalf("before step got %v want %v", got, want)
			}
			i.Next()
		}
	})

	t.Run("repeat-two-way", func(t *testing.T) {
		i := NewInterval(On(0), On(2), On(1), RepeatFromToFrom)
		for _, want := range []int{0, 1, 2, 1, 0, 1, 2} {
			if got := i.Value(); got != want {
				t.Fatalf("before step got %v want %v", got, want)
			}
			i.Next()
		}
	})
}

func TestInterval_Value_Storex(t *testing.T) {
	i := NewInterval(On(0), On(1), On(-1), RepeatFromTo)
	if got, want := i.Storex(), "interval(0,1,-1,'repeat')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInterval_Inspect(t *testing.T) {
	i := NewInterval(On(0), On(4), On(2), RepeatFromTo)
	in := Inspection{Properties: map[string]any{}}
	if got, want := i.Value(), 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := i.Storex(), "interval(0,4,2,'repeat')"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	i.Inspect(in)
	if got, want := in.Properties["value"], 0; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := in.Properties["length"], 3; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestInterval_StrategyHelpers(t *testing.T) {
	tests := []struct {
		name     string
		strategy string
		id       int
	}{
		{name: "once", strategy: "once", id: OnceFromTo},
		{name: "repeat", strategy: "repeat", id: RepeatFromTo},
		{name: "two-way", strategy: "two-way", id: OnceFromToFrom},
		{name: "repeat-two-way", strategy: "repeat-two-way", id: RepeatFromToFrom},
		{name: "once", strategy: "unknown", id: OnceFromTo},
	}

	for _, tc := range tests {
		if got := ParseIntervalStrategy(tc.strategy).id(); got != tc.id {
			t.Fatalf("ParseIntervalStrategy(%q) = %d, want %d", tc.strategy, got, tc.id)
		}
		if got := intervalStrategyName(tc.id); got != tc.name {
			t.Fatalf("intervalStrategyName(%d) = %q, want %q", tc.id, got, tc.name)
		}
		if got := asIntervalStrategy(tc.id).id(); got != tc.id {
			t.Fatalf("asIntervalStrategy(%d) = %d, want %d", tc.id, got, tc.id)
		}
	}

	if got, want := intervalStrategyName(999), "?"; got != want {
		t.Fatalf("intervalStrategyName(999) = %q, want %q", got, want)
	}
}
