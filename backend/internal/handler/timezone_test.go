package handler

import (
	"testing"
	"time"
)

// Dates without an offset must be read in the server's timezone (TZ),
// including across DST changes.
func TestParseDateTimeUsesLocalTimezone(t *testing.T) {
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Skip("tzdata not available")
	}
	orig := time.Local
	time.Local = london
	defer func() { time.Local = orig }()

	cases := []struct {
		in    string
		isEnd bool
		want  string
	}{
		{"2026-07-01", false, "2026-06-30T23:00:00Z"}, // BST start of day
		{"2026-07-01", true, "2026-07-01T22:59:59Z"},  // BST end of day
		{"2026-01-15", true, "2026-01-15T23:59:59Z"},  // GMT end of day
		{"2026-03-29", true, "2026-03-29T22:59:59Z"},  // 23-hour day (clocks forward)
		{"2026-10-25", true, "2026-10-25T23:59:59Z"},  // 25-hour day (clocks back)
		{"2026-07-01T12:00:00", false, "2026-07-01T11:00:00Z"},
	}
	for _, c := range cases {
		got := parseDateTime(c.in, c.isEnd)
		if got == nil {
			t.Fatalf("parseDateTime(%q) returned nil", c.in)
		}
		if s := got.UTC().Format(time.RFC3339); s != c.want {
			t.Errorf("parseDateTime(%q, %v) = %s, want %s", c.in, c.isEnd, s, c.want)
		}
	}

	start, end := parseTimeRange("2026-03-29", "2026-03-29", 24)
	if s := start.Format(time.RFC3339); s != "2026-03-29T00:00:00Z" {
		t.Errorf("parseTimeRange start = %s", s)
	}
	if s := end.Format(time.RFC3339); s != "2026-03-29T22:59:59Z" {
		t.Errorf("parseTimeRange end = %s", s)
	}
}
