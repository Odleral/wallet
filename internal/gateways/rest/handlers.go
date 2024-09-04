package rest

func (s *Server) endpoints() {
	s.router.GET("ping", s.pong())
	//s.router.GET("/balance/:id", s.balance())

	//s.router.POST("/transfer", s.transfer())
	//s.router.GET("/transaction/:id", s.transactionByID())

	s.router.GET("/exists/:id", s.exists())
}
