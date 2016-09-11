package giggleconf

type GiggleRconf struct {
}

func NewGiggleRconf() *GiggleRconf {
  return &GiggleRconf{
  }
}

func ReadGiggleRConf() *GiggleRconf {
  return NewGiggleRconf()
}
