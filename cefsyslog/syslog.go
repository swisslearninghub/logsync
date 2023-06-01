// Copyright 2023 Swiss Learning Hub AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cefsyslog

// Simplified replacement for default log/syslog component to meet requirements

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Priority int

const severityMask = 0x07
const facilityMask = 0xf8

const (
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

const (
	LOG_KERN Priority = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

// A Writer is a connection to a remote syslog server.
type Writer struct {
	priority Priority
	tag      string
	hostname string
	network  string
	raddr    string
	stamp    time.Time
	mu       sync.Mutex // guards conn
	conn     net.Conn
}

// SyslogWriterDial establishes connection to remote log daemon
func SyslogWriterDial(network, raddr string, priority Priority, tag string) (*Writer, error) {
	if network == "" {
		return nil, errors.New("local logging not implemented")
	}
	if priority < 0 || priority > LOG_LOCAL7|LOG_DEBUG {
		return nil, errors.New("invalid priority")
	}

	if tag == "" {
		tag = os.Args[0]
	}
	hostname, _ := os.Hostname()

	w := &Writer{
		priority: priority,
		tag:      tag,
		hostname: hostname,
		network:  network,
		raddr:    raddr,
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	err := w.connect()
	if err != nil {
		return nil, err
	}
	return w, err
}

// Write sends a log message to the syslog daemon.
func (w *Writer) Write(b []byte) (int, error) {
	return w.writeAndRetry(w.priority, string(b))
}

// Close closes a connection to the syslog daemon.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		err := w.conn.Close()
		w.conn = nil
		return err
	}
	return nil
}

// Log writes log using given timestamp
func (w *Writer) Log(stamp time.Time, p Priority, m string) error {
	w.stamp = stamp
	_, err := w.writeAndRetry(p, m)
	return err
}

// connect makes a connection to the syslog server. It must be called with w.mu held.
func (w *Writer) connect() (err error) {
	if w.conn != nil {
		// ignore err from close, it makes sense to continue anyway
		_ = w.conn.Close()
		w.conn = nil
	}

	var c net.Conn
	c, err = net.Dial(w.network, w.raddr)
	if err == nil {
		w.conn = c
		if w.hostname == "" {
			w.hostname = c.LocalAddr().String()
		}
	}

	return
}

func (w *Writer) writeAndRetry(p Priority, s string) (int, error) {
	pr := (w.priority & facilityMask) | (p & severityMask)

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.conn != nil {
		if n, err := w.write(pr, s); err == nil {
			return n, err
		}
	}
	if err := w.connect(); err != nil {
		return 0, err
	}
	return w.write(pr, s)
}

// write generates and writes a syslog formatted string. The
// format is as follows: <PRI>TIMESTAMP HOSTNAME TAG[PID]: MSG
func (w *Writer) write(p Priority, msg string) (int, error) {
	// ensure it ends in a \n
	nl := ""
	if !strings.HasSuffix(msg, "\n") {
		nl = "\n"
	}
	// pass on timestamp configured in writer
	err := w.writeString(w.stamp, p, w.hostname, w.tag, msg, nl)
	if err != nil {
		return 0, err
	}
	// Note: return the length of the input, not the number of
	// bytes printed by Fprintf, because this must behave like
	// an io.Writer.
	return len(msg), nil
}

func (w *Writer) writeString(stamp time.Time, p Priority, hostname, tag, msg, nl string) error {
	_, err := fmt.Fprintf(w.conn, "<%d>%s %s %s[%d]: %s%s",
		p, stamp.Format(time.RFC3339), hostname,
		tag, os.Getpid(), msg, nl)
	return err
}
