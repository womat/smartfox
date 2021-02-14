package smartfox

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type RelaySession struct {
	client *Client
	number int
}

type RelayStatus struct {
	Mode    uint16
	Control uint16
}

const MaxRelayNr = 4

const (
	On     = 1
	Off    = 0
	Auto   = 1
	Manual = 2
)

var ErrInvalidRelay = errors.New("smartfox: invalid relay number")

func (c *Client) NewRelay(n int) (*RelaySession, error) {
	if n > MaxRelayNr || n < 1 {
		return nil, ErrInvalidRelay
	}

	return &RelaySession{client: c, number: n}, nil
}

// relay mode 0=Off, 1=Auto, 2=Man. On
func (r *RelaySession) Status() (rs RelayStatus, err error) {
	if rs.Control, err = r.client.getRelayControl(r.number); err != nil {
		return rs, err
	}
	if rs.Mode, err = r.client.getRelayMode(r.number); err != nil {
		return rs, err
	}
	return rs, nil
}

func (r *RelaySession) SetControlOff() (RelayStatus, error) {
	if err := r.client.setRelayControl(r.number, Off); err != nil {
		return RelayStatus{}, err
	}
	return r.Status()
}

func (r *RelaySession) SetControlOn() (RelayStatus, error) {
	if err := r.client.setRelayControl(r.number, On); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

// Relay x Module
func (r *RelaySession) On() (RelayStatus, error) {
	if err := r.client.setRelayMode(r.number, Manual); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

func (r *RelaySession) Off() (RelayStatus, error) {
	if err := r.client.setRelayMode(r.number, Off); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

func (r *RelaySession) Auto() (RelayStatus, error) {
	if err := r.client.setRelayMode(r.number, Auto); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

// Relay x Module
func (c *Client) getRelayMode(nr int) (uint16, error) {
	return c.getRelay(nr, "RELAY_%d_MODE")
}

func (c *Client) getRelayControl(nr int) (uint16, error) {
	return c.getRelay(nr, "CONTROL_RELAY_%d")
}

func (c *Client) setRelayControl(nr int, ctrl uint16) error {
	return c.setRelay(nr, ctrl, "CONTROL_RELAY_%d")
}

func (c *Client) setRelayMode(nr int, ctrl uint16) error {
	return c.setRelay(nr, ctrl, "RELAY_%d_MODE")
}

func (c *Client) getRelay(nr int, register string) (uint16, error) {
	if nr > MaxRelayNr || nr < 1 {
		return 0, ErrInvalidRelay
	}

	b, err := c.readHoldingRegisters(fmt.Sprintf(register, nr))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(b), nil
}

func (c *Client) setRelay(nr int, ctrl uint16, register string) error {
	if nr > MaxRelayNr || nr < 1 {
		return ErrInvalidRelay
	}

	data := make([]uint8, 2)
	binary.BigEndian.PutUint16(data, ctrl)

	return c.writeHoldingRegisters(fmt.Sprintf(register, nr), data)
}
