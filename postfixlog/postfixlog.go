package postfixlog

import (
	"fmt"
	"regexp"
	"strconv"
)

type Stats struct {
	Delays            float64
	ReceivingDelay    float64
	QueuingDelay      float64
	ConnectionDelay   float64
	TransmissionDelay float64
	DSN               int
}

func bFloat64(b []byte) float64 {
	f, _ := strconv.ParseFloat(string(b), 64)
	return f
}

func bInt(b []byte) int {
	i, _ := strconv.Atoi(string(b))
	return i
}

// Apr 19 12:50:52 relaymail1 postfix/smtp[7570]: 69FFFC00B6: to=<xxxxxxx@example.jp>, relay=x.x.x.x[y.y.y.y]:25, delay=0.31, delays=0.04/0/0.09/0.17, dsn=2.0.0, status=sent (250 Ok)

var re = regexp.MustCompile(`, delay=(.+?), delays=(.+?)/(.+?)/(.+?)/(.+?), dsn=(\d)\.`)

// Parse :
func Parse(d1 []byte) (*Stats, error) {
	rs := re.FindSubmatch(d1)
	if len(rs) == 0 {
		return nil, fmt.Errorf("Not matched")
	}
	return &Stats{
		bFloat64(rs[1]),
		bFloat64(rs[2]),
		bFloat64(rs[3]),
		bFloat64(rs[4]),
		bFloat64(rs[5]),
		bInt(rs[6]),
	}, nil
}
