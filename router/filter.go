package router

func FilterChain(filters ...FilterFunc) FilterFunc {
	return func(next Handle) Handle {
		for i := len(filters) - 1; i >= 0; i-- {
			next = filters[i](next)
		}
		return next
	}
}
