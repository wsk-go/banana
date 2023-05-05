package defines

type BeanLoaded interface {
	Loaded()
}

type BeanConfiguration interface {
	Configuration() ModuleFunc
}
