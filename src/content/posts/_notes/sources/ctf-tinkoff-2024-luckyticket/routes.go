func (r *Router) initRoutes() {
	// API Router Group
	apiGroup := r.router.Group("/api")

	unauthorizedApiGroup := apiGroup.Group("")

	authorizedApiGroup := apiGroup.Group("")
	authorizedApiGroup.Use(middleware.AuthRequired)

	// Auth controllers
	{
		unauthorizedApiGroup.Handle(http.MethodPost, "/auth/register", api.RegisterUser(r.store))
		unauthorizedApiGroup.Handle(http.MethodPost, "/auth/login", api.LoginUser(r.store))
		authorizedApiGroup.Handle(http.MethodPost, "/auth/logout", api.LogoutUser())
	}

	// User controllers
	{
		authorizedApiGroup.Handle(http.MethodGet, "/user/me", api.GetUserSelfInfo(r.store))
		authorizedApiGroup.Handle(http.MethodGet, "/user/tickets", api.GetUserSelfTickets(r.store))
	}

	// Ticket controllers
	{
		authorizedApiGroup.Handle(http.MethodGet, "/ticket/:id", api.FindTicketByID(r.store))
		authorizedApiGroup.Handle(http.MethodPost, "/ticket/buy", api.BuyTicket(r.store))
		authorizedApiGroup.Handle(http.MethodPost, "/ticket/eat", api.EatTicket(r.store))
	}

	// Interview controllers
	{
		authorizedApiGroup.Handle(http.MethodPost, "/interview", api.CheckLuckForInterview(r.store))
	}
}
