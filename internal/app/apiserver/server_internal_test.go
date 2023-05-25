package apiserver

import (
	"bytes"
	"employee-service/internal/app/store/teststore"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Passport struct {
	Type   string
	Number string
}

type Department struct {
	Name  string
	Phone string
}

func TestServer_HandleEmployeeCreate(t *testing.T) {
	s := newServer(teststore.New())

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"name":       "Андрей",
				"surname":    "Андреев",
				"phone":      "+123456",
				"companyID":  1,
				"passport":   Passport{Type: "multi", Number: "4442"},
				"department": Department{Name: "Отдел №45", Phone: "99-99"},
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		// TODO валидация параметров, проверка на их отсутсвие
		{
			name: "invalid params",
			payload: map[string]interface{}{
				"surname":    "Андреев",
				"phone":      "+123456",
				"companyID":  1,
				"passport":   Passport{Type: "multi", Number: "4442"},
				"department": Department{Name: "Отдел №45", Phone: "99-99"},
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/employee", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
