package secondclass

type Fact struct {
	Value  string `yaml:"value"`
	Name   string `yaml:"name"`
	Unique string `yaml:"unique"`
}

func NewFact(name string, value string) *Fact {
	return &Fact{
		Name:   name,
		Value:  value,
		Unique: name + value,
	}
}
