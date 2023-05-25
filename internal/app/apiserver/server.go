package apiserver

import (
	"employee-service/internal/app/model"
	"employee-service/internal/app/store"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/employee", s.handleEmployeeCreate()).Methods("POST")
	//s.router.HandleFunc("/company/{id:[0-9]+}/", s.handleEmployeesByCompany()).Methods("GET")
	//s.router.HandleFunc("/department/{id:[0-9]+}/", s.handleEmployeesByDepartment()).Methods("GET")
	s.router.HandleFunc("/employee/{id:[0-9]+}/", s.handlerEmployeeDelete()).Methods("DELETE")
	//s.router.HandleFunc("/employee/{id:[0-9]+}/", s.handleEmployeeUpdate()).Methods("PATCH")
}

func (s *server) handleEmployeeCreate() http.HandlerFunc {
	type request struct {
		Name      string `json:"name"`
		Surname   string `json:"surname"`
		Phone     string `json:"phone"`
		CompanyID int    `json:"companyID"`
		Passport  struct {
			Type   string
			Number string
		} `json:"passport"`
		Department struct {
			Name  string
			Phone string
		} `json:"department"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		e := &model.Employee{
			Name:      req.Name,
			Surname:   req.Surname,
			Phone:     req.Phone,
			CompanyID: req.CompanyID,
			Passport: struct {
				Type   string
				Number string
			}{Type: req.Passport.Type, Number: req.Passport.Number},
			Department: struct {
				Name  string
				Phone string
			}{Name: req.Department.Name, Phone: req.Department.Phone},
		}

		if err := s.store.Employee().Create(e); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, fmt.Sprintf("id: %d", e.ID))
	}
}

func (s *server) handlerEmployeeDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handling delete employee at %s\n", r.URL.Path)

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		if err := s.store.Employee().Delete(id); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

//func (es *employeeServer) lastNameHandler(w http.ResponseWriter, req *http.Request) {
//	log.Printf("handling employee by lastName at %s\n", req.URL.Path)
//
//	lastName := mux.Vars(req)["lastName"]
//	employees, err := es.storage.GetEmployeesByLastName(lastName)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	renderJSON(w, employees)
//}

//func (es *employeeServer) updateEmployeeHandler(w http.ResponseWriter, req *http.Request) {
//	log.Printf("handling employee update at %s\n", req.URL.Path)
//
//	// Types used internally in this handler to (de-)serialize the request and
//	// response from/to JSON.
//	type RequestEmployee struct {
//		FirstName string `json:"firstName"`
//		LastName  string `json:"lastName"`
//		Email     string `json:"email"`
//	}
//
//	// Enforce a JSON Content-Type.
//	contentType := req.Header.Get("Content-Type")
//	mediatype, _, err := mime.ParseMediaType(contentType)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	if mediatype != "application/json" {
//		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
//		return
//	}
//
//	dec := json.NewDecoder(req.Body)
//	dec.DisallowUnknownFields()
//	var re RequestEmployee
//	if err := dec.Decode(&re); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	id, _ := strconv.Atoi(mux.Vars(req)["id"])
//	err = es.storage.UpdateEmployee(id, re.FirstName, re.LastName, re.Email)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
