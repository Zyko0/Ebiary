package core

type errorBase struct {
	msg string
}

func (e *errorBase) Error() string {
	return e.msg
}

type panicError struct {
	*errorBase
}

type shaderError struct {
	*errorBase
	line int
	char int
}

type imageError struct {
	*errorBase
}

type shaderFileError struct {
	*errorBase
}