package core

const (
	bufferSize = 40
)

type buffer struct {
	sources []string
	index   int
}

func newBuffer() *buffer {
	return &buffer{
		sources: make([]string, 0, bufferSize),
		index:   0,
	}
}

func (b *buffer) Empty() bool {
	return len(b.sources) == 0
}

func (b *buffer) Cancel() string {
	b.index = max(b.index-1, 0)
	return b.sources[b.index]
}

func (b *buffer) Redo() string {
	b.index = min(b.index+1, len(b.sources)-1)
	return b.sources[b.index]
}

func (b *buffer) New(src string) {
	b.sources = b.sources[:min(b.index+1, len(b.sources))]
	if len(b.sources) < bufferSize {
		b.sources = append(b.sources, src)
	} else {
		b.sources = append(b.sources[1:], src)
	}

	b.index = len(b.sources) - 1
}
