package postfixlog

import "testing"

// Test parser for postfix log
func TestParse(t *testing.T) {
	log := []byte("Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<xxxxxxx@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.31, delays=0.04/0/0.09/0.17, dsn=2.0.0, status=sent (250 Ok)")
	stats, err := Parse(log)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if stats.DSN != 2 {
		t.Errorf("Expected DSN 2, got %d", stats.DSN)
	}
}
