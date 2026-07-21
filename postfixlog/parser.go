package postfixlog

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/montanaflynn/stats"
)

type StatsBin struct {
	delays            []float64
	receivingDelay    []float64
	queuingDelay      []float64
	connectionDelay   []float64
	transmissionDelay []float64
	c2xx              float64
	c4xx              float64
	c5xx              float64
	total             float64
	duration          float64
}

var logFilter = []byte(" postfix/smtp[")

type PostfixParser struct {
	StatsBin *StatsBin
}

func NewPostfixParser() *PostfixParser {
	return &PostfixParser{
		StatsBin: NewStatsBin(),
	}
}

func (pp *PostfixParser) Parse(b []byte) error {
	if bytes.Index(b, logFilter) < 0 {
		return nil
	}
	s, err := Parse(b)
	if err != nil {
		return err
	}
	pp.StatsBin.Append(s)
	return nil
}

func (pp *PostfixParser) Finish(duration float64) {
	pp.StatsBin.duration = duration
}

func (pp *PostfixParser) Output(nowTime time.Time) string {
	return pp.StatsBin.Output(nowTime)
}

func NewStatsBin() *StatsBin {
	return &StatsBin{
		delays:            []float64{},
		receivingDelay:    []float64{},
		queuingDelay:      []float64{},
		connectionDelay:   []float64{},
		transmissionDelay: []float64{},
		c2xx:              0,
		c4xx:              0,
		c5xx:              0,
		total:             0,
		duration:          0,
	}
}

func (sb *StatsBin) Append(s *Stats) {
	switch s.DSN {
	case 2:
		sb.c2xx++
	case 4:
		sb.c4xx++
	case 5:
		sb.c5xx++
	}
	sb.total++

	sb.delays = append(sb.delays, s.Delays)
	sb.receivingDelay = append(sb.receivingDelay, s.ReceivingDelay)
	sb.queuingDelay = append(sb.queuingDelay, s.QueuingDelay)
	sb.connectionDelay = append(sb.connectionDelay, s.ConnectionDelay)
	sb.transmissionDelay = append(sb.transmissionDelay, s.TransmissionDelay)
}

// DisplayDelay :
func (sb *StatsBin) OutputDelay(w io.Writer, key string, arr []float64, now uint64) {
	if len(arr) > 0 {
		mean, _ := stats.Mean(arr)
		fmt.Fprintf(w, "postfixlog.%s_delay.average\t%f\t%d\n", key, mean, now)
		p99, _ := stats.Percentile(arr, 99)
		p95, _ := stats.Percentile(arr, 95)
		p90, _ := stats.Percentile(arr, 90)
		fmt.Fprintf(w, "postfixlog.%s_delay.99_percentile\t%f\t%d\n", key, p99, now)
		fmt.Fprintf(w, "postfixlog.%s_delay.95_percentile\t%f\t%d\n", key, p95, now)
		fmt.Fprintf(w, "postfixlog.%s_delay.90_percentile\t%f\t%d\n", key, p90, now)
	}
}

// Display :
func (sb *StatsBin) Output(nowTime time.Time) string {
	var buf bytes.Buffer
	now := uint64(nowTime.Unix())
	sb.OutputDelay(&buf, "total", sb.delays, now)
	sb.OutputDelay(&buf, "recving", sb.receivingDelay, now)
	sb.OutputDelay(&buf, "queuing", sb.queuingDelay, now)
	sb.OutputDelay(&buf, "connection", sb.connectionDelay, now)
	sb.OutputDelay(&buf, "transmission", sb.transmissionDelay, now)

	if sb.duration > 0 {
		fmt.Fprintf(&buf, "postfixlog.transfer_num.2xx_count\t%f\t%d\n", sb.c2xx/sb.duration, now)
		fmt.Fprintf(&buf, "postfixlog.transfer_num.4xx_count\t%f\t%d\n", sb.c4xx/sb.duration, now)
		fmt.Fprintf(&buf, "postfixlog.transfer_num.5xx_count\t%f\t%d\n", sb.c5xx/sb.duration, now)
		fmt.Fprintf(&buf, "postfixlog.transfer_total.count\t%f\t%d\n", sb.total/sb.duration, now)
	}
	if sb.total > 0 {
		fmt.Fprintf(&buf, "postfixlog.transfer_ratio.2xx_percentage\t%f\t%d\n", sb.c2xx*100/sb.total, now)
		fmt.Fprintf(&buf, "postfixlog.transfer_ratio.4xx_percentage\t%f\t%d\n", sb.c4xx*100/sb.total, now)
		fmt.Fprintf(&buf, "postfixlog.transfer_ratio.5xx_percentage\t%f\t%d\n", sb.c5xx*100/sb.total, now)
	}
	return buf.String()
}
