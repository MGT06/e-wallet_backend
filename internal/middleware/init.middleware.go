package middleware

type Middleware struct {
	OriginsAllowed []string
}

func InitMiddleWare(origins []string) *Middleware {
	return &Middleware{
		OriginsAllowed: origins,
	}
}
