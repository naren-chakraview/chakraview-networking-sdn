package northbound

import (
	"fmt"
	"net/http"
	"sync"
)

type APIServer struct {
	mu     sync.RWMutex
	mux    *http.ServeMux
	routes map[string]http.Handler
}

func NewAPIServer() *APIServer {
	return &APIServer{
		mux:    http.NewServeMux(),
		routes: make(map[string]http.Handler),
	}
}

func (as *APIServer) RegisterHandler(path string, handler http.Handler) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.routes[path] = handler
	as.mux.Handle(path, handler)
	fmt.Printf("Registered API endpoint: %s\n", path)
}

func (as *APIServer) ListRoutes() []string {
	as.mu.RLock()
	defer as.mu.RUnlock()

	routes := make([]string, 0, len(as.routes))
	for path := range as.routes {
		routes = append(routes, path)
	}
	return routes
}

func (as *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	as.mux.ServeHTTP(w, r)
}
