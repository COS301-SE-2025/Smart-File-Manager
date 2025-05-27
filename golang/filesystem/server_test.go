package filesystem

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddDirectoryHandler(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(getCompositeHandler))
	defer ts.Close()

	// Test cases
	testCases := []struct {
		name         string
		id           string
		dirName      string
		path         string
		expectedCode int
	}{
		{
			name:         "Valid Request",
			id:           "testid",
			dirName:      "testname",
			path:         "../testRootFolder",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Missing Parameters",
			id:           "",
			dirName:      "",
			path:         "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("POST",
				ts.URL+"/addDirectory?id="+tc.id+
					"&name="+tc.dirName+
					"&path="+tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Make request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d",
					tc.expectedCode, resp.StatusCode)
			}
		})
	}
}

func TestRemoveDirectoryHandler(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(removeCompositeHandler))
	defer ts.Close()

	// Test cases
	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{
			name:         "Valid Request",
			path:         "../testRootFolder",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Missing Parameters",
			path:         "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("DELETE",
				ts.URL+"/removeDirectory?path="+tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Make request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// Check status code
			if resp.StatusCode != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d",
					tc.expectedCode, resp.StatusCode)
			}
		})
	}
}
