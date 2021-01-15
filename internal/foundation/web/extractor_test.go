package web_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/bmizerany/assert"

	. "github.com/texazcowboy/warehouse/internal/foundation/web"
)

func TestExtractIntFromRequest(t *testing.T) {
	tests := []struct {
		name            string
		pathVarName     string
		pathVarValue    string
		isErrorExpected bool
	}{
		{
			name:            "Success",
			pathVarName:     "tt",
			pathVarValue:    "123",
			isErrorExpected: false,
		},
		{
			name:            "Invalid variable type",
			pathVarName:     "tt",
			pathVarValue:    "string_value",
			isErrorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			givenVars := map[string]string{tt.pathVarName: tt.pathVarValue}
			givenRequest := mux.SetURLVars(httptest.NewRequest("", "/test/path/"+tt.pathVarValue, nil), givenVars)
			// when
			extractedPathVarValue, err := ExtractIntFromRequest(givenRequest, tt.pathVarName)
			// then
			if err != nil {
				if tt.isErrorExpected {
					return
				}
				t.Errorf("Test: %v Failed with: %v", tt.name, err)
				return
			}
			assert.Equal(t, tt.pathVarValue, fmt.Sprintf("%d", extractedPathVarValue))
		})
	}
}
