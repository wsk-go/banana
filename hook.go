package banana

type BeanLoaded interface {
	BeanLoaded(application *Banana) error
}

// ApplicationLoaded call after all bean has been registered and injected
type ApplicationLoaded interface {
	ApplicationLoaded(application *Banana) error
}
