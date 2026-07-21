package postfixlog

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

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

			// テスト結果の検証
			if statsBin.c2xx != tt.expected.c2xx {
				t.Errorf("c2xx = %v, want %v", statsBin.c2xx, tt.expected.c2xx)
			}
			if statsBin.c4xx != tt.expected.c4xx {
				t.Errorf("c4xx = %v, want %v", statsBin.c4xx, tt.expected.c4xx)
			}
			if statsBin.c5xx != tt.expected.c5xx {
				t.Errorf("c5xx = %v, want %v", statsBin.c5xx, tt.expected.c5xx)
			}
			if statsBin.total != tt.expected.total {
				t.Errorf("total = %v, want %v", statsBin.total, tt.expected.total)
			}

			// スライスの比較
			if len(statsBin.delays) != len(tt.expected.delays) {
				t.Errorf("delays length = %v, want %v", len(statsBin.delays), len(tt.expected.delays))
			} else {
				for i, v := range statsBin.delays {
					if v != tt.expected.delays[i] {
						t.Errorf("delays[%d] = %v, want %v", i, v, tt.expected.delays[i])
					}
				}
			}

			if len(statsBin.receivingDelay) != len(tt.expected.receivingDelay) {
				t.Errorf("receivingDelay length = %v, want %v", len(statsBin.receivingDelay), len(tt.expected.receivingDelay))
			} else {
				for i, v := range statsBin.receivingDelay {
					if v != tt.expected.receivingDelay[i] {
						t.Errorf("receivingDelay[%d] = %v, want %v", i, v, tt.expected.receivingDelay[i])
					}
				}
			}

			if len(statsBin.queuingDelay) != len(tt.expected.queuingDelay) {
				t.Errorf("queuingDelay length = %v, want %v", len(statsBin.queuingDelay), len(tt.expected.queuingDelay))
			} else {
				for i, v := range statsBin.queuingDelay {
					if v != tt.expected.queuingDelay[i] {
						t.Errorf("queuingDelay[%d] = %v, want %v", i, v, tt.expected.queuingDelay[i])
					}
				}
			}

			if len(statsBin.connectionDelay) != len(tt.expected.connectionDelay) {
				t.Errorf("connectionDelay length = %v, want %v", len(statsBin.connectionDelay), len(tt.expected.connectionDelay))
			} else {
				for i, v := range statsBin.connectionDelay {
					if v != tt.expected.connectionDelay[i] {
						t.Errorf("connectionDelay[%d] = %v, want %v", i, v, tt.expected.connectionDelay[i])
					}
				}
			}

			if len(statsBin.transmissionDelay) != len(tt.expected.transmissionDelay) {
				t.Errorf("transmissionDelay length = %v, want %v", len(statsBin.transmissionDelay), len(tt.expected.transmissionDelay))
			} else {
				for i, v := range statsBin.transmissionDelay {
					if v != tt.expected.transmissionDelay[i] {
						t.Errorf("transmissionDelay[%d] = %v, want %v", i, v, tt.expected.transmissionDelay[i])
					}
				}
			}
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

	// 各行が正しい形式であることを確認
	for _, line := range lines {
		if !strings.Contains(line, "postfixlog.total_delay.") {
			t.Errorf("Unexpected line format: %s", line)
		}
		if !strings.Contains(line, "\t") {
			t.Errorf("Line missing timestamp: %s", line)
		}
	}
}
