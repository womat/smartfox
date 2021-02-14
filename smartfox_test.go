package smartfox

import (
	"fmt"
	"testing"
	"time"
)

const connection = "TCP 192.168.65.197:502 device:1 timeout:2 retries:3"

func TestGetPerformance(t *testing.T) {
	c := New()
	if err := c.Connect(connection); err != nil {
		t.Errorf("cann't connect to smartfox:%v", err)
		return
	}
	p, err := c.GetPerformance()
	if err != nil {
		fmt.Printf("%#v", p)
		t.Errorf("cann't get performance data sunrise time:%v", err)
		return
	}

	c.Close()
}

func TestRelay(t *testing.T) {
	c := New()
	if err := c.Connect(connection); err != nil {
		t.Errorf("cann't connect to smartfox:%v", err)
		return
	}

	r, err := c.NewRelay(1)
	if err != nil {
		t.Errorf("cann't connect to Relay:%v", err)
		return
	}

	if _, err := r.SetControlOff(); err != nil {
		t.Errorf("cann't switch to auto:%v", err)
		return
	}

	if _, err := r.Off(); err != nil {
		t.Errorf("cann't switch off:%v", err)
		return
	}
	time.Sleep(time.Second)
	s, err := r.Status()
	if err != nil || s.Mode != Off || s.Control != Off {
		t.Errorf("cann't switch off:%v", err)
		return
	}

	if _, err = r.On(); err != nil {
		t.Errorf("cann't switch on:%v", err)
		return
	}
	time.Sleep(time.Second)
	s, err = r.Status()
	if err != nil || s.Mode != Manual {
		t.Errorf("cann't switch on:%v", err)
		return
	}

	if _, err := r.Off(); err != nil {
		t.Errorf("cann't switch off:%v", err)
		return
	}
	time.Sleep(time.Second)
	s, err = r.Status()
	if err != nil || s.Mode != Off {
		t.Errorf("cann't switch off:%v", err)
		return
	}

	if _, err := r.Auto(); err != nil {
		t.Errorf("cann't switch to auto:%v", err)
		return
	}

	c.Close()
}
