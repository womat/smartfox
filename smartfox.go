package smartfox

import (
	"errors"
	"fmt"
	"time"

	"github.com/goburrow/modbus"
	"github.com/womat/tools"
)

const (
	On     = 1
	Off    = 0
	Auto   = 1
	Manual = 2
)

var (
	ErrInvalidLength   = errors.New("smartfox: invalid result length")
	ErrUnknownRegister = errors.New("smartfox: unknown register")
	ErrInvalidResponse = errors.New("smartfox: invalid modbus response")
	ErrInvalidRelay    = errors.New("smartfox: invalid relay number")
)

// Client structure contains all Properties of a connection
type Client struct {
	connectionString string
	timeout          time.Duration
	deviceID         uint8
	maxRetries       int
	clientHandler    *modbus.TCPClientHandler
	mbClient         modbus.Client
}

type Performance struct {
	Info struct {
		SwVersion             int
		ModbusProtocolVersion int
		ControlViaModbus      int
		Wlan                  struct {
			SwVersion  int
			MacAddress string
			Rssi       int
		}
		Lan struct {
			MacAddress string
		}
	}
	Frequency float64
	Power     struct {
		Total                        float64
		Smartfox                     float64
		L1, L2, L3                   float64
		FactorL1, FactorL2, FactorL3 float64
	}
	Voltage struct {
		L1, L2, L3 float64
	}
	Current struct {
		L1, L2, L3 float64
	}
	Energy struct {
		Smartfox                            float64
		ToGrid, FromGrid                    float64
		DayToGrid, DayFromGrid, DaySmartfox float64
	}
	Converter [5]struct {
		Power, Energy float64
	}
	Analog struct {
		AoutMode       int     // Aout mode
		Aout           float64 // Aout Output %
		ControlOutU    float64 // Control Analog out U %
		ControlOutI    float64 // Control Analog out I %
		ControlVoltage float64 // Control Voltage 24V
	}
}

func New() *Client {
	return &Client{}
}

func (c *Client) String() string {
	return "modbus"
}

// "TCP 192.168.65.197:502 device:1 timeout:2 retries:3
func (c *Client) Connect(connectionString string) (err error) {
	tools.GetField(&c.connectionString, connectionString, "connection")
	tools.GetField(&c.deviceID, connectionString, "device")
	tools.GetField(&c.timeout, connectionString, "timeout")
	//TODO: change to retry?
	tools.GetField(&c.maxRetries, connectionString, "retries")

	c.clientHandler = modbus.NewTCPClientHandler(c.connectionString)
	c.clientHandler.SlaveId = c.deviceID

	if err = c.clientHandler.Connect(); err != nil {
		return
	}

	c.mbClient = modbus.NewClient(c.clientHandler)
	return
}

func (c *Client) ReConnect() error {
	if err := c.clientHandler.Connect(); err != nil {
		return err
	}

	c.mbClient = modbus.NewClient(c.clientHandler)
	return nil
}

func (c *Client) Close() error {
	return c.clientHandler.Close()
}

// Session Performance
func (c *Client) GetPerformance() (p Performance, err error) {
	if p.Info.SwVersion, err = c.readInt("SW_VERSION"); err != nil {
		return
	}
	if p.Info.ModbusProtocolVersion, err = c.readInt("MODBUS_VERSION"); err != nil {
		return
	}
	if p.Info.ControlViaModbus, err = c.readInt("CONTROL_MODBUS"); err != nil {
		return
	}
	if p.Info.Wlan.SwVersion, err = c.readInt("WLAN_VERSION"); err != nil {
		return
	}
	if p.Info.Wlan.MacAddress, err = c.WLanMacAddress(); err != nil {
		return
	}
	if p.Info.Wlan.Rssi, err = c.readInt("WIFI_RSSI"); err != nil {
		return
	}
	if p.Info.Lan.MacAddress, err = c.LanMacAddress(); err != nil {
		return
	}

	if p.Frequency, err = c.readFloat64("FREQUENCY"); err != nil {
		return
	}

	if p.Power.Smartfox, err = c.readFloat64("POWER_SMARTFOX"); err != nil {
		return
	}
	if p.Power.Total, err = c.readFloat64("POWER_TOTAL"); err != nil {
		return
	}
	if p.Power.L1, err = c.readFloat64("POWER_L1"); err != nil {
		return
	}
	if p.Power.L2, err = c.readFloat64("POWER_L2"); err != nil {
		return
	}
	if p.Power.L3, err = c.readFloat64("POWER_L3"); err != nil {
		return
	}

	if p.Power.FactorL1, err = c.readFloat64("POWERFACTOR_L1"); err != nil {
		return
	}
	if p.Power.FactorL2, err = c.readFloat64("POWERFACTOR_L2"); err != nil {
		return
	}
	if p.Power.FactorL3, err = c.readFloat64("POWERFACTOR_L3"); err != nil {
		return
	}

	if p.Voltage.L1, err = c.readFloat64("VOLTAGE_L1"); err != nil {
		return
	}
	if p.Voltage.L2, err = c.readFloat64("VOLTAGE_L2"); err != nil {
		return
	}
	if p.Voltage.L3, err = c.readFloat64("VOLTAGE_L3"); err != nil {
		return
	}

	if p.Current.L1, err = c.readFloat64("CURRENT_L1"); err != nil {
		return
	}
	if p.Current.L2, err = c.readFloat64("CURRENT_L2"); err != nil {
		return
	}
	if p.Current.L3, err = c.readFloat64("CURRENT_L3"); err != nil {
		return
	}

	if p.Energy.FromGrid, err = c.readFloat64("ENERGY_FROM_GRID"); err != nil {
		return
	}
	if p.Energy.ToGrid, err = c.readFloat64("ENERGY_TO_GRID"); err != nil {
		return
	}
	if p.Energy.Smartfox, err = c.readFloat64("ENERGY_SMARTFOX"); err != nil {
		return
	}

	if p.Energy.DayFromGrid, err = c.readFloat64("DAY_FROM_GRID"); err != nil {
		return
	}
	if p.Energy.DayToGrid, err = c.readFloat64("DAY_TO_GRID"); err != nil {
		return
	}
	if p.Energy.DaySmartfox, err = c.readFloat64("DAY_SMARTFOX"); err != nil {
		return
	}

	for i := range p.Converter {
		if p.Converter[i].Energy, err = c.readFloat64(fmt.Sprintf("CONVERTER_%d_ENERGY", i+1)); err != nil {
			return
		}
		if p.Converter[i].Power, err = c.readFloat64(fmt.Sprintf("CONVERTER_%d_POWER", i+1)); err != nil {
			return
		}
	}

	if p.Analog.AoutMode, err = c.readInt("AOUT_MODE"); err != nil {
		return
	}
	if p.Analog.Aout, err = c.readFloat64("AOUT"); err != nil {
		return
	}
	if p.Analog.ControlOutI, err = c.readFloat64("CONTROL_ANALOG_I"); err != nil {
		return
	}
	if p.Analog.ControlOutU, err = c.readFloat64("CONTROL_ANALOG_U"); err != nil {
		return
	}
	if p.Analog.ControlVoltage, err = c.readFloat64("CONTROL_24V"); err != nil {
		return
	}

	return
}

// Version
func (c *Client) Version() (i int, err error) {
	return c.readInt("SW_VERSION")
}
