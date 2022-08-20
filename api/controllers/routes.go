package controllers

import "github.com/jumkos/WartaGkjMedari/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Renungan routes
	s.Router.HandleFunc("/renungan", middlewares.SetMiddlewareJSON(s.CreateRenungan)).Methods("POST")
	s.Router.HandleFunc("/renungan", middlewares.SetMiddlewareJSON(s.GetAllRenungan)).Methods("GET")
	s.Router.HandleFunc("/renungan/{id}", middlewares.SetMiddlewareJSON(s.GetRenungan)).Methods("GET")
	s.Router.HandleFunc("/renungan/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateRenungan))).Methods("PUT")
	s.Router.HandleFunc("/renungan/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteRenungan)).Methods("DELETE")
}