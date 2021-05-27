package smartfox

import (
	"encoding/binary"
	"fmt"
)

const maxRelayNr = 4

type Relay struct {
	Client *Client
	number int
}

type RelayStatus struct {
	Mode    uint16
	Control uint16
}

func (c *Client) NewRelay(n int) (*Relay, error) {
	if n > maxRelayNr || n < 1 {
		return nil, ErrInvalidRelay
	}

	return &Relay{Client: c, number: n}, nil
}

// Status returns the relay state
//  mode 0=Off, 1=Auto, 2=Man. On
func (r *Relay) Status() (rs RelayStatus, err error) {
	if rs.Control, err = r.Client.getRelayControl(r.number); err != nil {
		return rs, err
	}
	if rs.Mode, err = r.Client.getRelayMode(r.number); err != nil {
		return rs, err
	}
	return rs, err
}

func (r *Relay) ControlOff() (RelayStatus, error) {
	if err := r.Client.setRelayControl(r.number, Off); err != nil {
		return RelayStatus{}, err
	}
	return r.Status()
}

func (r *Relay) ControlOn() (RelayStatus, error) {
	if err := r.Client.setRelayControl(r.number, On); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

func (r *Relay) On() (RelayStatus, error) {
	if err := r.Client.setRelayMode(r.number, Manual); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

func (r *Relay) Off() (RelayStatus, error) {
	if err := r.Client.setRelayMode(r.number, Off); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

func (r *Relay) Auto() (RelayStatus, error) {
	if err := r.Client.setRelayMode(r.number, Auto); err != nil {
		return RelayStatus{}, err
	}

	return r.Status()
}

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
	if nr > maxRelayNr || nr < 1 {
		return 0, ErrInvalidRelay
	}

	b, err := c.readHoldingRegisters(fmt.Sprintf(register, nr))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(b), nil
}

func (c *Client) setRelay(nr int, ctrl uint16, register string) error {
	if nr > maxRelayNr || nr < 1 {
		return ErrInvalidRelay
	}

	data := make([]uint8, 2)
	binary.BigEndian.PutUint16(data, ctrl)

	return c.writeHoldingRegisters(fmt.Sprintf(register, nr), data)
}
