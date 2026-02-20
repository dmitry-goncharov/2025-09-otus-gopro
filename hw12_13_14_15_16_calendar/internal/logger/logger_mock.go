package logger

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (l *Mock) Debug(_ string) {}
func (l *Mock) Info(_ string)  {}
func (l *Mock) Warn(_ string)  {}
func (l *Mock) Error(_ string) {}
