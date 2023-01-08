package httpstatus

type Type int

const (
	Unknown     Type = 0
	Information Type = 1
	Success     Type = 2
	Redirection Type = 3
	ClientError Type = 4
	ServerError Type = 5
)

func From(code int) Type {
	switch code / 100 {
	case 1:
		return Information
	case 2:
		return Success
	case 3:
		return Redirection
	case 4:
		return ClientError
	case 5:
		return ServerError
	default:
		return Unknown
	}
}
