package limiter

type FakeClock struct {
	currentTime int64
}

func (f *FakeClock) Now() int64 {
	return f.currentTime
}

func (f *FakeClock) Advance(seconds int64) {
	f.currentTime += seconds
}
