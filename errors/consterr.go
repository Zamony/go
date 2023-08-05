package errors

type ConstantError string

func (e ConstantError) Error() string {
	return string(e)
}

func (e ConstantError) Is(target error) bool {
	switch target := target.(type) {
	case ConstantError:
		return target == e
	default:
		return false
	}
}

func (e ConstantError) As(target any) bool {
	switch target := target.(type) {
	case *ConstantError:
		*target = e
		return true
	default:
		return false
	}
}
