package errors

type ConstantError string

func (e ConstantError) Error() string { return string(e) }
