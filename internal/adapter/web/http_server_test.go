package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leroxyl/computer-manager-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockComputerManager struct {
	mock.Mock
}

func (m *MockComputerManager) Create(computer entity.Computer) error {
	args := m.MethodCalled("Create", computer)
	return args.Error(0)
}

func (m *MockComputerManager) Read(mac string) (entity.Computer, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockComputerManager) Update(computer entity.Computer) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockComputerManager) Delete(mac string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockComputerManager) ReadAll() ([]entity.Computer, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockComputerManager) ReadAllForEmployee(abbr string) ([]entity.Computer, error) {
	//TODO implement me
	panic("implement me")
}

func TestServer_createComputer(t *testing.T) {
	testCases := []struct {
		name           string
		setExpectation func(m *MockComputerManager)
		body           []byte
		respCode       int
		respBody       string
	}{
		{
			name: "happy path",
			setExpectation: func(m *MockComputerManager) {
				m.On("Create", entity.Computer{
					MACAddr:      "00:1B:44:11:3A:B7",
					ComputerName: "lenovo",
					IPAddr:       "192.158.1.38",
					EmployeeAbbr: "mmu",
					Description:  "text",
				}).Return(nil)
			},
			body:     []byte("{\"macAddr\": \"00:1B:44:11:3A:B7\", \"computerName\": \"lenovo\", \"ipAddr\": \"192.158.1.38\", \"employeeAbbr\": \"mmu\", \"description\": \"text\"}"),
			respCode: http.StatusOK,
			respBody: "{\"macAddr\":\"00:1B:44:11:3A:B7\",\"computerName\":\"lenovo\",\"ipAddr\":\"192.158.1.38\",\"employeeAbbr\":\"mmu\",\"description\":\"text\"}",
		},
		{
			name:           "invalid json",
			setExpectation: func(m *MockComputerManager) {},
			body:           []byte("{"),
			respCode:       http.StatusBadRequest,
			respBody:       "{\"error\":\"unexpected EOF\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCM := &MockComputerManager{}
			tc.setExpectation(mockCM)

			server := NewServer(mockCM)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/computers", bytes.NewBuffer(tc.body))

			server.router.ServeHTTP(w, req)

			assert.Equal(t, tc.respCode, w.Code)
			assert.Equal(t, tc.respBody, w.Body.String())
		})
	}
}
