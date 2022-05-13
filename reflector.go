package reflectricity

type Reflector struct {
	private            bool
	packages           map[string]string
	arrayMergeStrategy arrayMergeStrategy
}

func NewReflector(private bool) *Reflector {
	return &Reflector{
		private:            private,
		packages:           make(map[string]string),
		arrayMergeStrategy: CONCAT,
	}
}

func (r *Reflector) ArrayStrategy(ams arrayMergeStrategy) {
	r.arrayMergeStrategy = ams
}
