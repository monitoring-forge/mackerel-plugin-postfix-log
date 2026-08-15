package postfixlog

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// assertEqual checks if two float64 values are equal.
func assertEqual(t *testing.T, got, want float64, field string) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", field, got, want)
	}
}

// assertSliceEqual checks if two float64 slices are equal.
func assertSliceEqual(t *testing.T, got, want []float64, name string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s length = %d, want %d", name, len(got), len(want))
		return
	}
	for i, v := range got {
		if v != want[i] {
			t.Errorf("%s[%d] = %v, want %v", name, i, v, want[i])
		}
	}
}

// assertStatsBinEqual verifies all fields of StatsBin.
func assertStatsBinEqual(t *testing.T, statsBin, expected *StatsBin) {
	t.Helper()
	assertEqual(t, statsBin.c2xx, expected.c2xx, "c2xx")
	assertEqual(t, statsBin.c4xx, expected.c4xx, "c4xx")
	assertEqual(t, statsBin.c5xx, expected.c5xx, "c5xx")
	assertEqual(t, statsBin.total, expected.total, "total")
	assertSliceEqual(t, statsBin.delays, expected.delays, "delays")
	assertSliceEqual(t, statsBin.receivingDelay, expected.receivingDelay, "receivingDelay")
	assertSliceEqual(t, statsBin.queuingDelay, expected.queuingDelay, "queuingDelay")
	assertSliceEqual(t, statsBin.connectionDelay, expected.connectionDelay, "connectionDelay")
	assertSliceEqual(t, statsBin.transmissionDelay, expected.transmissionDelay, "transmissionDelay")
}

func TestPostfixParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *StatsBin
	}{
		{
			name:  "正常なログエントリ",
			input: []byte("Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<xxxxxxx@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.31, delays=0.04/0/0.09/0.17, dsn=2.0.0, status=sent (250 Ok)"),
			expected: &StatsBin{
				c2xx:              1,
				total:             1,
				delays:            []float64{0.31},
				receivingDelay:    []float64{0.04},
				queuingDelay:      []float64{0},
				connectionDelay:   []float64{0.09},
				transmissionDelay: []float64{0.17},
			},
		},
		{
			name:  "4xxエラーのログエントリ",
			input: []byte("Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<xxxxxxx@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.45, delays=0.05/0.1/0.15/0.15, dsn=4.0.0, status=deferred (450 Mailbox unavailable)"),
			expected: &StatsBin{
				c4xx:              1,
				total:             1,
				delays:            []float64{0.45},
				receivingDelay:    []float64{0.05},
				queuingDelay:      []float64{0.1},
				connectionDelay:   []float64{0.15},
				transmissionDelay: []float64{0.15},
			},
		},
		{
			name:  "5xxエラーのログエントリ",
			input: []byte("Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<xxxxxxx@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.55, delays=0.06/0.12/0.18/0.19, dsn=5.0.0, status=bounced (550 User unknown)"),
			expected: &StatsBin{
				c5xx:              1,
				total:             1,
				delays:            []float64{0.55},
				receivingDelay:    []float64{0.06},
				queuingDelay:      []float64{0.12},
				connectionDelay:   []float64{0.18},
				transmissionDelay: []float64{0.19},
			},
		},
		{
			name:     "フィルタに一致しないログエントリ",
			input:    []byte("Apr 19 12:50:52 relaymail1 postfix/qmgr[7570]: 69FFFC00B6: removed"),
			expected: &StatsBin{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statsBin := NewStatsBin()
			parser := &PostfixParser{StatsBin: statsBin}

			err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			assertStatsBinEqual(t, statsBin, tt.expected)
		})
	}
}

func TestNewStatsBin(t *testing.T) {
	sb := NewStatsBin()
	if sb == nil {
		t.Fatal("NewStatsBin returned nil")
	}
	if len(sb.delays) != 0 {
		t.Errorf("delays = %v, want empty slice", sb.delays)
	}
	if sb.c2xx != 0 || sb.c4xx != 0 || sb.c5xx != 0 || sb.total != 0 || sb.duration != 0 {
		t.Errorf("initial values = %+v, want all zeros", sb)
	}
}

func TestStatsBin_Append(t *testing.T) {
	tests := []struct {
		name          string
		stats         *Stats
		expectedC2xx  float64
		expectedC4xx  float64
		expectedC5xx  float64
		expectedTotal float64
	}{
		{
			name:          "DSN 2xx",
			stats:         &Stats{Delays: 0.31, DSN: 2},
			expectedC2xx:  1,
			expectedTotal: 1,
		},
		{
			name:          "DSN 4xx",
			stats:         &Stats{Delays: 0.45, DSN: 4},
			expectedC4xx:  1,
			expectedTotal: 1,
		},
		{
			name:          "DSN 5xx",
			stats:         &Stats{Delays: 0.55, DSN: 5},
			expectedC5xx:  1,
			expectedTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStatsBin()
			sb.Append(tt.stats)

			assertEqual(t, sb.c2xx, tt.expectedC2xx, "c2xx")
			assertEqual(t, sb.c4xx, tt.expectedC4xx, "c4xx")
			assertEqual(t, sb.c5xx, tt.expectedC5xx, "c5xx")
			assertEqual(t, sb.total, tt.expectedTotal, "total")
		})
	}
}

func TestStatsBin_OutputDelay(t *testing.T) {
	statsBin := NewStatsBin()
	statsBin.delays = []float64{0.1, 0.2, 0.3, 0.4, 0.5}

	var buf bytes.Buffer
	now := uint64(time.Now().Unix())
	statsBin.OutputDelay(&buf, "total", statsBin.delays, now)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 4 {
		t.Errorf("Expected 4 lines of output, got %d", len(lines))
	}

	for _, line := range lines {
		if !strings.Contains(line, "postfixlog.total_delay.") {
			t.Errorf("Unexpected line format: %s", line)
		}
		if !strings.Contains(line, "\t") {
			t.Errorf("Line missing timestamp: %s", line)
		}
	}
}

func TestStatsBin_OutputDelay_EmptySlice(t *testing.T) {
	sb := NewStatsBin()
	var buf bytes.Buffer
	now := uint64(time.Now().Unix())
	sb.OutputDelay(&buf, "total", sb.delays, now)

	if buf.Len() != 0 {
		t.Errorf("Expected empty output for empty slice, got %q", buf.String())
	}
}

func TestStatsBin_Output(t *testing.T) {
	sb := NewStatsBin()
	sb.c2xx = 80
	sb.c4xx = 15
	sb.c5xx = 5
	sb.total = 100
	sb.duration = 60
	sb.delays = []float64{0.1, 0.2, 0.3}
	sb.receivingDelay = []float64{0.05}
	sb.queuingDelay = []float64{0.1}
	sb.connectionDelay = []float64{0.15}
	sb.transmissionDelay = []float64{0.2}

	output := sb.Output(time.Now())

	expectedKeys := []string{
		"postfixlog.total_delay.average",
		"postfixlog.total_delay.99_percentile",
		"postfixlog.total_delay.95_percentile",
		"postfixlog.total_delay.90_percentile",
		"postfixlog.transfer_num.2xx_count",
		"postfixlog.transfer_num.4xx_count",
		"postfixlog.transfer_num.5xx_count",
		"postfixlog.transfer_total.count",
		"postfixlog.transfer_ratio.2xx_percentage",
		"postfixlog.transfer_ratio.4xx_percentage",
		"postfixlog.transfer_ratio.5xx_percentage",
	}

	for _, key := range expectedKeys {
		if !strings.Contains(output, key) {
			t.Errorf("Output missing expected key: %s", key)
		}
	}
}

func TestPostfixParser_Output(t *testing.T) {
	pp := NewPostfixParser()
	log := "Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<test@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.31, delays=0.04/0/0.09/0.17, dsn=2.0.0, status=sent (250 Ok)"

	err := pp.Parse([]byte(log))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	pp.StatsBin.duration = 60
	output := pp.Output(time.Now())

	if output == "" {
		t.Error("Output should not be empty")
	}
}

func TestPostfixParser_Finish(t *testing.T) {
	pp := NewPostfixParser()
	pp.Finish(120)

	if pp.StatsBin.duration != 120 {
		t.Errorf("duration = %v, want 120", pp.StatsBin.duration)
	}
}

func BenchmarkPostfixParser_Parse(b *testing.B) {
	log := []byte("Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<test@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.31, delays=0.04/0/0.09/0.17, dsn=2.0.0, status=sent (250 Ok)")
	pp := NewPostfixParser()

	for b.Loop() {
		_ = pp.Parse(log)
	}
}
