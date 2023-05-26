package apiserver

import (
	"employee-service/internal/app/model"
	"employee-service/internal/app/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"mime"
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
	s.router.HandleFunc("/company/{id:[0-9]+}/", s.handleEmployeesByCompany()).Methods("GET")
	s.router.HandleFunc("/company/{id:[0-9]+}/department/{name}/", s.handleEmployeesByDepartment()).Methods("GET")
	s.router.HandleFunc("/employee/{id:[0-9]+}/", s.handlerEmployeeDelete()).Methods("DELETE")
	s.router.HandleFunc("/employee/{id:[0-9]+}/", s.handleEmployeeUpdate()).Methods("PATCH")
}

// handleEmployeeCreate создает запись о новом сотруднике в БД и возвращает его ID
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

		s.respond(w, r, http.StatusCreated, &e.ID)
	}
}

// handleEmployeeCreate удаляет сотрудника из БД по его ID
func (s *server) handlerEmployeeDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("handling delete employee at %s\n", r.URL.Path)

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		if err := s.store.Employee().Delete(id); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusNoContent, nil)
	}
}

// handleEmployeesByCompany возвращает записи о всех сотрудниках по указанному в запросе ID компании
func (s *server) handleEmployeesByCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("handling employee by companyID at %s\n", r.URL.Path)

		companyID, _ := strconv.Atoi(mux.Vars(r)["id"])

		employees, err := s.store.Employee().FindByCompany(companyID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if len(employees) == 0 {
			s.respond(w, r, http.StatusNotFound, nil)
			return
		}

		s.respond(w, r, http.StatusOK, employees)
	}
}

// handleEmployeesByDepartment возвращает записи о всех сотрудниках по указанному в запросе ID компании и названию отдела
func (s *server) handleEmployeesByDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("handling employee by companyID at %s\n", r.URL.Path)

		companyID, _ := strconv.Atoi(mux.Vars(r)["id"])
		department := mux.Vars(r)["name"]

		employees, err := s.store.Employee().FindByDepartment(companyID, department)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if len(employees) == 0 {
			s.respond(w, r, http.StatusNotFound, nil)
			return
		}

		s.respond(w, r, http.StatusOK, employees)
	}
}

// handleEmployeeUpdate принимает в запросе новые данные сотрудника и обновляет запись в БД по соотвествующему ID
func (s *server) handleEmployeeUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("handling employee update at %s\n", r.URL.Path)

		// Enforce a JSON Content-Type.
		contentType := r.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if mediatype != "application/json" {
			http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		re := new(model.Employee)
		if err = dec.Decode(re); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		presentData, err := s.store.Employee().FindById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if re.Name == "" {
			re.Name = presentData.Name
		}
		if re.Surname == "" {
			re.Surname = presentData.Surname
		}
		if re.Phone == "" {
			re.Phone = presentData.Phone
		}
		if re.CompanyID == 0 {
			re.CompanyID = presentData.CompanyID
		}
		if re.Passport.Type == "" {
			re.Passport.Type = presentData.Passport.Type
		}
		if re.Passport.Number == "" {
			re.Passport.Number = presentData.Passport.Number
		}
		if re.Department.Name == "" {
			re.Department.Name = presentData.Department.Name
		}
		if re.Department.Phone == "" {
			re.Department.Phone = presentData.Department.Phone
		}

		if err = re.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.store.Employee().PartiallyUpdate(id, re)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		s.respond(w, r, http.StatusOK, re)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
