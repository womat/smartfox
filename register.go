package smartfox

import (
	"encoding/binary"
	"fmt"
)

type register struct {
	address  uint16
	typ      string
	quantity uint16
	scale    float64
}

var registerMap = map[string]register{
	"SW_VERSION":         {40010, "uint32", 2, 1},
	"WLAN_VERSION":       {40012, "uint32", 2, 1},
	"MAC_LAN":            {40014, "uint8", 3, 1},
	"MODBUS_VERSION":     {40017, "uint16", 1, 1},
	"MAC_WIFI":           {40018, "uint8", 3, 1},
	"CONTROL_MODBUS":     {40400, "uint8", 1, 1},
	"CONTROL_ANALOG_U":   {40401, "uint8", 1, 0.1},
	"CONTROL_ANALOG_I":   {40402, "uint8", 1, 0.1},
	"CONTROL_RELAY_1":    {40403, "uint8", 1, 1},
	"CONTROL_RELAY_2":    {40404, "uint8", 1, 1},
	"CONTROL_RELAY_3":    {40405, "uint8", 1, 1},
	"CONTROL_RELAY_4":    {40406, "uint8", 1, 1},
	"ENERGY_FROM_GRID":   {41000, "uint64", 4, 1},
	"ENERGY_TO_GRID":     {41004, "uint64", 4, 1},
	"ENERGY_SMARTFOX":    {41008, "uint64", 4, 1},
	"DAY_ENERGY":         {41012, "uint32", 6, 1},
	"DAY_FROM_GRID":      {41012, "uint32", 2, 1},
	"DAY_TO_GRID":        {41014, "uint32", 2, 1},
	"DAY_SMARTFOX":       {41016, "uint32", 2, 1},
	"POWER_TOTAL":        {41018, "int32", 2, 1},
	"POWER_L1":           {41020, "int32", 2, 1},
	"POWER_L2":           {41022, "int32", 2, 1},
	"POWER_L3":           {41024, "int32", 2, 1},
	"VOLTAGE_L1":         {41026, "uint16", 1, 0.1},
	"VOLTAGE_L2":         {41027, "uint16", 1, 0.1},
	"VOLTAGE_L3":         {41028, "uint16", 1, 0.1},
	"CURRENT_L1":         {41029, "uint32", 2, 0.001},
	"CURRENT_L2":         {41031, "uint32", 2, 0.001},
	"CURRENT_L3":         {41033, "uint32", 2, 0.001},
	"POWERFACTOR_L1":     {41035, "int16", 1, 0.0001},
	"POWERFACTOR_L2":     {41036, "int16", 1, 0.0001},
	"POWERFACTOR_L3":     {41037, "int16", 1, 0.0001},
	"FREQUENCY":          {41038, "uint16", 1, 0.01},
	"PT1000":             {41039, "int16", 1, 0.1},
	"S0_INPUT":           {41040, "uint32", 2, 1},
	"POWER_SMARTFOX":     {41042, "uint32", 2, 1},
	"CONTROL_24V":        {41044, "uint16", 1, 0.01},
	"WIFI_RSSI":          {41046, "int8", 1, 1},
	"AOUT":               {41047, "uint8", 1, 1},
	"CONVERTER_1_POWER":  {41400, "uint32", 2, 1},
	"CONVERTER_1_ENERGY": {41402, "uint64", 4, 1},
	"CONVERTER_2_POWER":  {41406, "uint32", 2, 1},
	"CONVERTER_2_ENERGY": {41408, "uint64", 4, 1},
	"CONVERTER_3_POWER":  {41412, "uint32", 2, 1},
	"CONVERTER_3_ENERGY": {41414, "uint64", 4, 1},
	"CONVERTER_4_POWER":  {41418, "uint32", 2, 1},
	"CONVERTER_4_ENERGY": {41420, "uint64", 4, 1},
	"CONVERTER_5_POWER":  {41424, "uint32", 2, 1},
	"CONVERTER_5_ENERGY": {41426, "uint64", 4, 1},
	"AOUT_MODE":          {42207, "uint8", 1, 1},
	"RELAY_1_MODE":       {42250, "uint8", 1, 1},
	"RELAY_2_MODE":       {42280, "uint8", 1, 1},
	"RELAY_3_MODE":       {42310, "uint8", 1, 1},
	"RELAY_4_MODE":       {42340, "uint8", 1, 1},
}

func (c *Client) readFloat64(key string) (float64, error) {
	b, err := c.readHoldingRegisters(key)
	if err != nil {
		return 0, err
	}

	register, ok := registerMap[key]
	if !ok {
		return 0, ErrUnknownRegister
	}

	switch register.typ {
	case "int8":
		return float64(int8(binary.BigEndian.Uint16(b))) * register.scale, nil
	case "uint8":
		return float64(binary.BigEndian.Uint16(b)) * register.scale, nil
	case "int16":
		return float64(int16(binary.BigEndian.Uint16(b))) * register.scale, nil
	case "uint16":
		return float64(binary.BigEndian.Uint16(b)) * register.scale, nil
	case "int32":
		return float64(int32(binary.BigEndian.Uint32(b))) * register.scale, nil
	case "uint32":
		return float64(binary.BigEndian.Uint32(b)) * register.scale, nil
	case "int64":
		return float64(int64(binary.BigEndian.Uint64(b))) * register.scale, nil
	case "uint64":
		return float64(binary.BigEndian.Uint64(b)) * register.scale, nil
	}

	return 0, ErrUnknownRegister
}

func (c *Client) readInt(key string) (int, error) {
	b, err := c.readHoldingRegisters(key)
	if err != nil {
		return 0, err
	}

	register, ok := registerMap[key]
	if !ok {
		return 0, ErrUnknownRegister
	}

	switch register.typ {
	case "int8":
		return int(float64(int8(binary.BigEndian.Uint16(b))) * register.scale), nil
	case "uint8":
		return int(float64(binary.BigEndian.Uint16(b)) * register.scale), nil
	case "int16":
		return int(float64(int16(binary.BigEndian.Uint16(b))) * register.scale), nil
	case "uint16":
		return int(float64(binary.BigEndian.Uint16(b)) * register.scale), nil
	case "int32":
		return int(float64(int32(binary.BigEndian.Uint32(b))) * register.scale), nil
	case "uint32":
		return int(float64(binary.BigEndian.Uint32(b)) * register.scale), nil
	case "int64":
		return int(float64(int64(binary.BigEndian.Uint64(b))) * register.scale), nil
	case "uint64":
		return int(float64(binary.BigEndian.Uint64(b)) * register.scale), nil
	}

	return 0, ErrUnknownRegister
}

// MAC-Address LAN
func (c *Client) LanMacAddress() (string, error) {
	b, err := c.readHoldingRegisters("MAC_LAN")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%02X-%02X-%02X-%02X-%02X-%02X", b[0], b[1], b[2], b[3], b[4], b[5]), nil
}

// MAC-Address Wifi
func (c *Client) WLanMacAddress() (string, error) {
	b, err := c.readHoldingRegisters("MAC_WIFI")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%02X-%02X-%02X-%02X-%02X-%02X", b[0], b[1], b[2], b[3], b[4], b[5]), nil
}

// Control AoutMode
func (c *Client) getAoutMode() (byte, error) {
	return c.readByte("AOUT_MODE")
}

func (c *Client) readByte(key string) (byte, error) {
	b, err := c.readHoldingRegisters(key)
	if err != nil {
		return 0, err
	}

	return byte(binary.BigEndian.Uint16(b)), nil
}

func (c *Client) readHoldingRegisters(key string) ([]uint8, error) {
	register, ok := registerMap[key]
	if !ok {
		return nil, ErrUnknownRegister
	}

	switch b, err := c.mbClient.ReadHoldingRegisters(register.address-1, register.quantity); {
	case err != nil:
		return b, err
	case len(b) != int(register.quantity)*2:
		return b, ErrInvalidLength
	default:
		return b, nil
	}
}

func (c *Client) writeHoldingRegisters(key string, regdata []uint8) error {
	register, ok := registerMap[key]
	if !ok {
		return ErrUnknownRegister
	}
	if len(regdata) != int(register.quantity)*2 {
		return ErrInvalidLength
	}

	switch b, err := c.mbClient.WriteMultipleRegisters(register.address-1, register.quantity, regdata); {
	case err != nil:
		return err
	case binary.BigEndian.Uint16(b) != register.quantity:
		return ErrInvalidResponse
	default:
		return nil
	}
}
