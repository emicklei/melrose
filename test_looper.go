package melrose

type TestLooper struct {
	Beats int64
	Bars  int64
	Biab  int64
}

func (t *TestLooper) Tick() {
	t.Beats++
	if t.Beats%t.Biab == 0 {
		t.Bars++
	}
}

func (t *TestLooper) Start() {}
func (t *TestLooper) Stop()  {}
func (t *TestLooper) Reset() {
	t.Beats = 0
	t.Bars = 0
}

func (t *TestLooper) SetBPM(bpm float64) {}
func (t *TestLooper) BPM() float64       { return 120.0 }

func (t *TestLooper) SetBIAB(biab int) {
	t.Biab = int64(biab)
}
func (t *TestLooper) BIAB() int {
	return int(t.Biab)
}

func (t *TestLooper) Begin(l *Loop) {}
func (t *TestLooper) End(l *Loop)   {}

func (t *TestLooper) BeatsAndBars() (int64, int64) {
	return t.Beats, t.Bars
}

func (t *TestLooper) Delay(bars int64, beats int64, seq Sequenceable) {

}
