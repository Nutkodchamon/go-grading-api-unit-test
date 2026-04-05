package grade

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-grading-api/config"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCalculateGrade_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		homework float64
		midterm  float64
		final    float64
		expected string
	}{
		{"Grade A", 80, 70, 90, "A"},
		{"Grade B", 70, 70, 70, "B"},
		{"Grade C", 60, 60, 60, "C"},
		{"Grade D", 50, 50, 50, "D"},
		{"Grade F", 40, 40, 40, "F"},
		{"Boundary A", 80, 80, 80, "A"},
		{"Invalid Negative HW", -1, 50, 50, "Invalid"},
		{"Invalid Negative Mid", 50, -1, 50, "Invalid"},
		{"Invalid Negative Final", 50, 50, -1, "Invalid"},
		{"Invalid Over 100 HW", 101, 50, 50, "Invalid"},
		{"Invalid Over 100 Mid", 50, 101, 50, "Invalid"},
		{"Invalid Over 100 Final", 50, 50, 101, "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, grade := CalculateGrade(tt.homework, tt.midterm, tt.final)
			assert.Equal(t, tt.expected, grade)
		})
	}
}

type MockRepository struct{}

func (m *MockRepository) GetGradeByStudentID(studentID string) (*Response, error) {
	if studentID == "db_error" {
		return nil, errors.New("database error")
	}
	return &Response{StudentID: studentID, Total: 85, Grade: "A"}, nil
}

func (m *MockRepository) InsertGrade(g Response, homework, midterm, final float64) error {
	if g.StudentID == "db_error" {
		return errors.New("database error")
	}
	return nil
}

type MockService struct{}

func (m *MockService) CheckGrade(studentID string) (*Response, error) {
	if studentID == "error" {
		return nil, errors.New("grade not found")
	}
	return &Response{StudentID: studentID, Total: 90, Grade: "A"}, nil
}

func (m *MockService) SubmitGrade(req Request) (*Response, error) {
	if req.StudentID == "error" {
		return nil, errors.New("failed to save grade")
	}
	return &Response{StudentID: req.StudentID, Total: 85, Grade: "A"}, nil
}

func TestNewGradeService(t *testing.T) {
	repo := &MockRepository{}
	s := NewGradeService(repo)
	assert.NotNil(t, s)
}

func TestSubmitGrade(t *testing.T) {
	repo := &MockRepository{}
	s := NewGradeService(repo)

	res, err := s.SubmitGrade(Request{StudentID: "123", Homework: 10, Midterm: 10, Final: 10})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res2, err2 := s.SubmitGrade(Request{StudentID: "db_error", Homework: 10, Midterm: 10, Final: 10})
	assert.Error(t, err2)
	assert.Nil(t, res2)
}

func TestCheckGrade(t *testing.T) {
	repo := &MockRepository{}
	s := NewGradeService(repo)

	res, err := s.CheckGrade("123")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res2, err2 := s.CheckGrade("")
	assert.Error(t, err2)
	assert.Nil(t, res2)

	res3, err3 := s.CheckGrade("db_error")
	assert.Error(t, err3)
	assert.Nil(t, res3)
}

func TestNewHandler(t *testing.T) {
	h := NewHandler(&MockService{})
	assert.NotNil(t, h)
}

func TestSubmitGradeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(&MockService{})
	router := gin.Default()
	router.POST("/grade/submit", h.SubmitGradeHandler)

	reqBody := bytes.NewBufferString(`{"studentId":"123","homework":10,"midterm":10,"final":10}`)
	req, _ := http.NewRequest("POST", "/grade/submit", reqBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	reqBodyBad := bytes.NewBufferString(`invalid_json_format`)
	reqBad, _ := http.NewRequest("POST", "/grade/submit", reqBodyBad)
	wBad := httptest.NewRecorder()
	router.ServeHTTP(wBad, reqBad)
	assert.Equal(t, http.StatusBadRequest, wBad.Code)

	reqBodyErr := bytes.NewBufferString(`{"studentId":"error","homework":10,"midterm":10,"final":10}`)
	reqErr, _ := http.NewRequest("POST", "/grade/submit", reqBodyErr)
	wErr := httptest.NewRecorder()
	router.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)
}

func TestGetGradeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(&MockService{})
	router := gin.Default()
	router.GET("/grade/:studentId", h.GetGradeHandler)

	req, _ := http.NewRequest("GET", "/grade/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	reqErr, _ := http.NewRequest("GET", "/grade/error", nil)
	wErr := httptest.NewRecorder()
	router.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusNotFound, wErr.Code)
}

func setupTestDB() {
	config.DB, _ = sql.Open("sqlite3", ":memory:")
	query := `
	CREATE TABLE IF NOT EXISTS grades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id TEXT,
		homework REAL,
		midterm REAL,
		final REAL,
		total REAL,
		grade TEXT
	);`
	config.DB.Exec(query)
}

func TestGradeRepository(t *testing.T) {
	setupTestDB()
	repo := &GradeRepository{}

	err := repo.InsertGrade(Response{StudentID: "123", Total: 100, Grade: "A"}, 100, 100, 100)
	assert.NoError(t, err)

	res, err := repo.GetGradeByStudentID("123")
	assert.NoError(t, err)
	assert.Equal(t, "123", res.StudentID)

	res2, err2 := repo.GetGradeByStudentID("999")
	assert.Error(t, err2)
	assert.Nil(t, res2)
}
