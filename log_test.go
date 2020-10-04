package crzerolog

import (
	"testing"
)

func TestTraceContextFromHeader(t *testing.T) {
	for _, tt := range []struct {
		header      string
		wantTraceID string
		wantSpanID  string
		wantSampled bool
	}{
		{"0123456789abcdef0123456789abcdef/123;o=1", "0123456789abcdef0123456789abcdef", "123", true},
		{"0123456789abcdef0123456789abcdef/123;o=0", "0123456789abcdef0123456789abcdef", "123", false},
		{"0123456789abcdef0123456789abcdef/123", "0123456789abcdef0123456789abcdef", "123", false},
		{"0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef", "", false},
		{"0123456789abcdef0123456789abcdef/invalid", "", "", false},
		{"invalid", "", "", false},
		{"", "", "", false},
	} {
		traceID, spanID, sampled := traceContextFromHeader(tt.header)
		if traceID != tt.wantTraceID || spanID != tt.wantSpanID || sampled != tt.wantSampled {
			t.Errorf("traceContextFromHeader(%q) = (%q, %q, %v), want = (%q, %q, %v)", tt.header, traceID, spanID, sampled, tt.wantTraceID, tt.wantSpanID, tt.wantSampled)
		}
	}
}
