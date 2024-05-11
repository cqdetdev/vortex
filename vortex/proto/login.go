package proto

type Login struct {
	Service string
	Token   string
}

func (l *Login) ID() uint32 {
	return LOGIN
}

func (l *Login) Marshal(io IO) {
	io.String(&l.Service)
	io.String(&l.Token)
}