package hook

type BeanLoaded interface {
	BeanLoaded() error
}

// ApplicationLoaded call after all bean has been registered and injected
type ApplicationLoaded interface {
	ApplicationLoaded() error
}
