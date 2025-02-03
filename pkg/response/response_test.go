package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		status     int
		wantStatus int
	}{
		{
			name:       "successful response",
			data:       map[string]string{"message": "success"},
			status:     http.StatusOK,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			JSON(rec, tt.data, tt.status)

			// Check status code
			if rec.Code != tt.wantStatus {
				t.Errorf("JSON() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			// Check Content-Type header
			contentType := rec.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("JSON() Content-Type = %v, want application/json", contentType)
			}

			// Verify JSON response
			var response map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Errorf("JSON() failed to decode response: %v", err)
			}
		})
	}
}
