package errors

type emailNotConfigured struct {
	msg string
}

func (e *emailNotConfigured) Error() string {
	return "email not configured: " + e.msg
}

func NewEmailNotConfigured(msg string) error {
	return &emailNotConfigured{msg: msg}
}

func IsEmailNotConfigured(err error) bool {
	_, ok := err.(*emailNotConfigured)
	return ok
}
