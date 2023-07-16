package option

type OptPrinter interface {
	String() string
}

type Option interface {
	Apply(OptPrinter)
}

type Opt struct {
	f func(OptPrinter)
}

func (fdo *Opt) Apply(do OptPrinter) {
	fdo.f(do)
}

func NewOpt(f func(OptPrinter)) *Opt {
	return &Opt{
		f: f,
	}
}
