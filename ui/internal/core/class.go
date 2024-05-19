package core

type ClassesImpl struct {
	classes []string
}

func (i *ClassesImpl) Classes() []string {
	return i.classes
}

func (i *ClassesImpl) SetClasses(classes ...string) {
	i.classes = append(i.classes[:0], classes...)
}

func (i *ClassesImpl) AddClasses(classes ...string) {
	i.classes = append(i.classes, classes...)
}
