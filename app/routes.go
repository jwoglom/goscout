package app

import "./endpoints"

func (s *Server) addRoutes() {
	s.router.HandleFunc("/status", endpoints.StatusHandler)
}
