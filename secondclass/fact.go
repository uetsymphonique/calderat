package secondclass

type Fact struct {
	Value  string `yaml:"value"`
	Trait  string `yaml:"trait"`
	Unique string `yaml:"unique"`
}

func NewFact(trait string, value string) *Fact {
	return &Fact{
		Trait:  trait,
		Value:  value,
		Unique: trait + value,
	}
}
