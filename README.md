กชมนต์ บริบูรณ์ 6609650145

---

# Go Grading API – Testing & Unit Testing Workshop

Welcome to the **Go Grading API Workshop**. This repository provides a hands-on environment to practice API integration testing and Go unit testing patterns.

---

## 🎯 Workshop Objectives
By the end of this session, you will be able to:
* Manage Go builds using a **Makefile**.
* Perform manual and automated API testing with **Postman**.
* Implement robust **Unit Tests** using Go's built-in testing package.

---

## 🚀 Getting Started

### Step 1: Fork the Repository
Navigate to the [original repository](https://github.com/go-training/go-grading-api-workshop) and click the **Fork** button to create your own copy.

### Step 2: Clone the Project
Replace `<your-username>` with your GitHub handle:
```bash
git clone [https://github.com/](https://github.com/)<your-username>/go-grading-api-unit-test
cd go-grading-api-unit-test
```

### Step 3: Build and Run
This project uses a Makefile to automate tasks. To compile and start the server, execute:
```bash
go run cmd/server/main.go
```
or
```bash
make run
```
---

## 🛠 Workshop API and Unit Test

Now that your server is running, follow these steps to test the Authentication and Grading logic.

---

### 🚀 Task 1: Manual API Testing (Postman)
You will verify the 4 core endpoints. Note that most endpoints require a **Bearer Token**.

1.  **Login & Get Token:**
    * Send the `POST Login` request with the default credentials (`John`/`1234`).
    * Copy the `token` (or `accessToken`) from the response body.
2.  **Authorize Other Requests:**
    * For the other 3 requests (Check token, Submit, Check Grade), go to the **Authorization** tab.
    * Select **Type: Bearer Token** and paste your token.
3.  **Test Grading Logic:**
    * **Submit:** Send `POST /api/grade/submit`. This should save the grade to your `university.db`.
    * **Check Grade:** Use the `GET /api/grade/:studentId` to verify the record exists.

---

### 🚀 Task 2: Postman Workflow & Automation
Manually pasting tokens is slow. Let's automate the "Login -> Store Token" workflow.

1.  **Automate Token Storage:**
    In the **Tests** tab of your **Login** request, add:
    ```javascript
    const response = pm.response.json();
    if (response.token) {
        pm.collectionVariables.set("jwt_token", response.token);
    }
    ```
2.  **Use Variable in Headers:**
    * Click on your **Collection** (grading-api-workshop) > **Authorization**.
    * Set Type to **Bearer Token**.
    * Set Token to `{{jwt_token}}`.
    * Ensure all individual requests are set to **"Inherit auth from parent"**.
3.  **Run Collection Runner:**
    Run the entire folder to ensure the sequence (Login -> Submit -> Check) works perfectly without manual intervention.

---

### 🧪 Task 3: Basic Unit Test
Since your project has a `pkg/jwt` and `internal/grade`, you should test the scoring and token logic.

1.  **Test Calculation Logic:**
    Create `internal/grade/grade_test.go` and test the scoring math:
    ```go
    package grade

    import "testing"

    func TestCalculateGradeA(t *testing.T) {
        _, grade := CalculateGrade(80, 70, 90)
    
        if grade != "A" {
            t.Error("Expected A")
        }
    }
    ```
    
2. **Run Tests:**
   Execute the following in your terminal:
    ```bash
    # Run a single test
    go test ./internal/grade -run TestCalculateGradeA
    
   # Run tests and save a coverage profile
    go test ./internal/grade -coverprofile=coverage.out

    # Open the visual report in your browser
    go tool cover -html=coverage.out
    ```
3. **Refactor to Use assert.Equal:**
   Refactor test to use `testify/assert` package:
   ```go
   package grade
   
   import (
       "testing"
       "github.com/stretchr/testify/assert"
   )
   
   func TestCalculateGradeA(t *testing.T) {
        _, grade := CalculateGrade(80, 70, 90)
        assert.Equal(t, "A", grade)
   }
   ```
   > **Note:** `go get github.com/stretchr/testify/assert` to install the package.

4.  **💡Expanding to 100% Code Coverage (20 minutes):**
    Create more test cases for the `CalculateGrade` function. 100% coverage is required.
    Expected test cases:
    ```text
    TestCalculateGradeA
    TestCalculateGradeB
    TestCalculateGradeC
    TestCalculateGradeD
    TestCalculateGradeF
    TestInvalidScore
    TestBoundaryScore
    ```
---

### 🧪 Task 4: Table-Driven Tests (TDD)
Refactor the `TestCalculateGradeA` test to use a table-driven approach:
```go
package grade

import (
	"testing"

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
		// TODO: add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, grade := CalculateGrade(tt.homework, tt.midterm, tt.final)
			assert.Equal(t, tt.expected, grade)
		})
	}
}
```
Run the test to ensure it passes:
```bash
make test
```

## Testing Styles Summary
| Testing Style        | Used For                          | Example                     |
|----------------------|-----------------------------------|-----------------------------|
| Basic Unit Test      | Test a single function            | TestCalculateGradeA         |
| Assert (testify)     | Make tests cleaner and readable   | assert.Equal()              |
| Table-Driven Test    | Test multiple input cases         | Multiple grade test cases   |
| Edge Case Test       | Test invalid or boundary inputs   | Negative score, >100        |

---

### 🧪 Task 5: Mock Testing (Handler & Service Tests)
We use mock objects to test handlers and services without calling the real database.

### 🧩 Service Test with Mock Database
1.  **Create Mock Database (For Service Test):**
    Create a mock repository `grade/mock_repository.go` to simulate database behavior.
```go
package grade

type MockRepository struct{}

func (m *MockRepository) GetGradeByStudentID(studentID string) (*Response, error) {
	return &Response{
		StudentID: studentID,
		Total:     85,
		Grade:     "A",
	}, nil
}

func (m *MockRepository) InsertGrade(g Response, homework, midterm, final float64) error {
	//TODO implement me
	panic("implement me")
}
```

2.  **Use mock repository in service test:**

```go
func TestCheckGrade(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewGradeService(mockRepo)

    res, err := service.CheckGrade("65001")

    assert.NoError(t, err)
    assert.Equal(t, "65001", res.StudentID)
    assert.Equal(t, "A", res.Grade)
}
```

Run the test:
```bash
go test ./internal/grade -run TestCheckGrade
```
> **Note:** This test does **not** connect to the database.

### 🧩 Handler Test with Mock Service
1. **Create Mock Service (For Handler Test):**
    Create a mock service `grade/mock_service.go` to simulate service behavior.
```go
package grade

type MockService struct{}

func (m *MockService) CheckGrade(studentID string) (*Response, error) {
	return &Response{
		StudentID: studentID,
		Total:     90,
		Grade:     "A",
	}, nil
}

func (m *MockService) SubmitGrade(req Request) (*Response, error) {
	//TODO implement me
	panic("implement me")
}
```
2. **Use mock service in handler test:**
    Create a mock handler `grade/mock_handler.go` to test the handler.
```go
package grade

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetGradeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockService{}
	handler := NewHandler(mockService)

	router := gin.Default()
	router.GET("/grade/:studentId", handler.GetGradeHandler)

	req, _ := http.NewRequest("GET", "/grade/65001", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

```

Run the test:
```bash
go test ./internal/grade -run TestGetGradeHandler
```

> **Note:** This test does **not** connect to the real service.