package controllers

import (
	// "context"
	// "encoding/json"
	// "io"
	// "io/ioutil"
	// "net/http"

	"github.com/baksman/rest-api/api/middleware"
	// "github.com/baksman/rest-api/api/responses"
)

// func middleware(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		ctx2 := context.WithValue(r.Context(), "userId", "3949949")
// 		handler.ServeHTTP(w, r.WithContext(ctx2))

// 		rByte, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			return
// 		}
// 		userName := struct{ Name string }{"baksman"}
// 		err = json.Unmarshal(rByte, &userName)
// 		if err != nil {
// 			responses.ERROR(w, 400, err)
// 			// log.Fatal(err)
// 		}
// 		io.WriteString(w, userName.Name)

// 	})
// }
func (s *Server) initializeRoutes() {
	// newRoute := s.Router.NewRoute()

	// Home Route
	s.Router.Use(middlewares.SetMiddlewareJSON)
	s.Router.HandleFunc("/", (s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", (s.Login)).Methods("POST")
	s.Router.HandleFunc("/sign-up", (s.Register)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", (s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", (s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", (s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", (middlewares.SetMiddlewareAuthentication(s.Updateuser))).Methods("PUT")
	s.Router.HandleFunc("/upload-dp", middlewares.SetMiddlewareAuthentication(s.uploadProfilePic)).Methods("DELETE")

	//Posts routes
	s.Router.HandleFunc("/posts", s.CreatePost).Methods("POST")
	s.Router.HandleFunc("/posts", (s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", (s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", (middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")
}
