package sub

import (
	"encoding/base64"
	"testing"

	"github.com/mhsanaei/3x-ui/v3/internal/database/model"
)

func TestApplyHappServerDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{name: "empty", description: "", want: ""},
		{name: "whitespace", description: "  \n ", want: ""},
		{name: "text", description: "  Описание сервера  ", want: base64.StdEncoding.EncodeToString([]byte("Описание сервера"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := map[string]string{}
			applyHappServerDescription(params, &model.Inbound{ServerDescription: tt.description})
			if got := params["serverDescription"]; got != tt.want {
				t.Fatalf("serverDescription = %q, want %q", got, tt.want)
			}
			if tt.want == "" {
				if _, exists := params["serverDescription"]; exists {
					t.Fatal("empty description must not add serverDescription")
				}
			}
		})
	}
}
