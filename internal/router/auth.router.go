package router

func (r *Router) AuthRouter() {
	auth := r.r.Group("auth")

	auth.POST("register", r.ac.RegisterAH)
}
