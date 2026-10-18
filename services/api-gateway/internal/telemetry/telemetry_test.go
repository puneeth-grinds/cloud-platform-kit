package telemetry

import (
	"reflect"
	"testing"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want map[string]string
	}{
		{
			name: "empty headers",
			raw:  "",
			want: map[string]string{},
		},
		{
			name: "one header",
			raw:  "Authorization=Basic token",
			want: map[string]string{"Authorization": "Basic token"},
		},
		{
			name: "multiple headers",
			raw:  "Authorization=Basic token,X-Tenant=123",
			want: map[string]string{
				"Authorization": "Basic token",
				"X-Tenant":      "123",
			},
		},
		{
			name: "value with equals sign",
			raw:  "Authorization=Basic abc=def",
			want: map[string]string{"Authorization": "Basic abc=def"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := parseHeaders(test.raw); !reflect.DeepEqual(got, test.want) {
				t.Errorf("parseHeaders(%q) = %#v, want %#v", test.raw, got, test.want)
			}
		})
	}
}
