package crzerolog

import (
	"testing"
)

func TestTraceContextFromHeader(t *testing.T) {
	for _, tt := range []struct {
		header      string
		wantTraceID string
		wantSpanID  string
	}{
		{"0123456789abcdef0123456789abcdef/123;o=1", "0123456789abcdef0123456789abcdef", "000000000000007b"},
		{"0123456789abcdef0123456789abcdef/123;o=0", "0123456789abcdef0123456789abcdef", "000000000000007b"},
		{"0123456789abcdef0123456789abcdef/123", "0123456789abcdef0123456789abcdef", "000000000000007b"},
		{"0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef", ""},
		{"0123456789abcdef0123456789abcdef/invalid", "", ""},
		{"invalid", "", ""},
		{"", "", ""},
	} {
		traceID, spanID := traceContextFromHeader(tt.header)
		if traceID != tt.wantTraceID || spanID != tt.wantSpanID {
			t.Errorf("traceContextFromHeader(%q) = (%q, %q), want = (%q, %q)", tt.header, traceID, spanID, tt.wantTraceID, tt.wantSpanID)
		}
	}
}
