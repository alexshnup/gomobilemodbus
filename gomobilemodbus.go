package gomobilemodbus

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goburrow/modbus"
	"github.com/shopspring/decimal"
)

var (
	DeviceSN string
)

var needDisableInputAction bool

func ModbusRequest(Serialport, SlaveId, Cmd, Address, Quantity string) (out string) {
	var n0, n1, n2 decimal.Decimal
	var err error

	needDisableInputAction = true

	if len(SlaveId) > 0 {
		n0, err = decimal.NewFromString(SlaveId)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
	}

	if len(Address) > 0 {
		n1, err = decimal.NewFromString(Address)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
	}

	if len(Quantity) > 0 {
		n2, err = decimal.NewFromString(Quantity)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
	}

	// Modbus RTU/ASCII
	// handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler := modbus.NewRTUClientHandler(Serialport)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 2
	handler.SlaveId = byte(n0.IntPart())
	handler.Timeout = 2 * time.Second

	err = handler.Connect()
	if err != nil {
		//log.Panicf("error %v", err)
		return "-1 " + fmt.Sprintf("%v", err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	//Check disabled inputs action
	if needDisableInputAction {
		// results, err := client.ReadHoldingRegisters(9, 1)
		results, err := client.ReadHoldingRegisters(8, 1)
		if err != nil {
			//log.Panicf("error ReadHoldingRegisters %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
		if results[1] == 0x00 {
			fmt.Println("need disable input actions. Autodisabling")
			results, err := client.WriteSingleRegister(9, 3)
			if err != nil {
				//log.Panicf("error %v", err)
				return "-1 " + fmt.Sprintf("%v", err)
			}
			fmt.Printf("%x", results[0])

			results, err = client.WriteSingleRegister(10, 3)
			if err != nil {
				//log.Panicf("error %v", err)
				return "-1 " + fmt.Sprintf("%v", err)
			}
			fmt.Printf("%x", results[0])

			results, err = client.WriteSingleRegister(8, 60)
			if err != nil {
				//log.Panicf("error %v", err)
				return "-1 " + fmt.Sprintf("%v", err)
			}
			fmt.Printf("%x", results[0])
		}

	}

	switch Cmd {

	case "rc":
		//go run main.go /dev/ttyUSB0 46 rc 0
		//0 or 1
		results, err := client.ReadCoils(uint16(n1.IntPart()), 1)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
		// fmt.Printf("%x", results[0])
		out = fmt.Sprintf("%x", results[0])

	case "rd":
		address := uint16(n1.IntPart())

		quantity := uint16(n2.IntPart())
		if quantity < 1 || quantity > 125 {
			quantity = 1
		}
		// go run main.go /dev/ttyUSB0 46 rd 0
		//0 or 1
		results, err := client.ReadDiscreteInputs(address, quantity)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
		if uint16(n1.IntPart()) < 2 {
			out = fmt.Sprintf("%x", results[0])
		} else {
			out = fmt.Sprintf("%x", results)
		}

	case "volt":
		//go run main.go /dev/ttyUSB0 46 volt
		//12.4V
		quantity := uint16(n2.IntPart())
		if quantity < 1 || quantity > 125 {
			quantity = 1
		}

		results, err := client.ReadInputRegisters(121, quantity)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}

		output, err := strconv.ParseInt(fmt.Sprintf("%x", results), 16, 64)
		if err != nil {
			//fmt.Println(err)
			return "-1 " + fmt.Sprintf("%v", err)
		}

		out = fmt.Sprintf("%.1fV", float32(float32(output)/1000))

	case "count":
		//go run main.go /dev/ttyUSB0 46 count 1
		//
		addr := uint16(n1.IntPart())
		if addr > 125 {
			addr = 0
		}

		results, err := client.ReadInputRegisters(32+addr, 1)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}

		output, err := strconv.ParseInt(fmt.Sprintf("%x", results), 16, 64)
		if err != nil {
			//fmt.Println(err)
			return "-1 " + fmt.Sprintf("%v", err)
		}

		out = fmt.Sprintf("%d", output)

	case "reset":
		//go run main.go /dev/ttyUSB0 46 reboot
		//

		client.WriteSingleRegister(120, 1)

	case "rh":
		//go run main.go /dev/ttyUSB0 46 rh 9
		//go run main.go /dev/ttyUSB0 46 rh 8

		address := uint16(n1.IntPart())

		quantity := uint16(n2.IntPart())
		if quantity < 1 || quantity > 125 {
			quantity = 1
		}

		// switch address {
		// case 16:
		// 	return
		// }

		//go run main.go /dev/ttyUSB0 46 rh 10 1
		//0003   if no   -  need write holding register
		results, err := client.ReadHoldingRegisters(address, quantity)
		if err != nil {
			//log.Panicf("error ReadHoldingRegisters %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}

		output, err := strconv.ParseInt(fmt.Sprintf("%x", results), 16, 64)
		if err != nil {
			//fmt.Println(err)
			return "-1 " + fmt.Sprintf("%v", err)
		}

		// fmt.Printf("0x%x", results)
		out = fmt.Sprintf("%d", output)

	case "wh":
		//go run main.go /dev/ttyUSB0 46 wh 9 0003
		results, err := client.WriteSingleRegister(uint16(n1.IntPart()), uint16(n2.IntPart()))
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
		fmt.Printf("%x", results[0])

	case "on":
		//go run main.go /dev/ttyUSB0 46 on 1
		//good responce is "255"
		results, err := client.WriteSingleCoil(uint16(n1.IntPart()), 0xFF00)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
		out = fmt.Sprintf("%v", results[0])

	case "off":
		//go run main.go /dev/ttyUSB0 46 off 1
		//0
		results, err := client.WriteSingleCoil(uint16(n1.IntPart()), 0x0000)
		if err != nil {
			//log.Panicf("error %v", err)
			return "-1 " + fmt.Sprintf("%v", err)
		}
		out = fmt.Sprintf("%v", results[0])
	}
	return out
}
