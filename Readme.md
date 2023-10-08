
# GoMobileModbus Library (Android)
## Introduction:
GoMobileModbus is a library designed to facilitate Modbus RTU/ASCII communication on Android devices using the Go programming language. It leverages existing packages for Modbus communication and extends functionality for a mobile environment.

some commands such as getting Voltage, designed for devise WB-MRM2-mini v.2.0
- Dual channel mini-module Relay
https://wirenboard.com/en/product/WB-MRM2-mini/

## Features:
- Simplified Modbus RTU/ASCII interaction with customizable parameters.
- Integration with the well-established github.com/goburrow/modbus package.
- Use of github.com/shopspring/decimal for precise decimal operations.
- Centralized error handling for better debugging and stability.

## Dependencies:
- github.com/goburrow/modbus: For core Modbus functionality.
- github.com/shopspring/decimal: For accurate decimal operations.

## Functions:
```go
ModbusRequest(Serialport, SlaveId, Cmd, Address, Quantity string) string
```
Handles all modbus operations. It takes the following parameters:

- Serialport: Path to the serial port (e.g., /dev/ttyUSB0).
- SlaveId: The ID of the slave device.
- Cmd: The command type to be executed (e.g., 'rc' for read coils).
- Address: Register/Coil address.
- Quantity: Number of Registers/Coils to read/write.

The function returns the output or an error message prefixed with "-1".


## Build:
```bash
go build -o gomobilemodbus.aar -target=android .
```
Don't use "go mod vendor" because it will create a vendor folder with all the dependencies and the gomobile tool will not be able to find the dependencies.

## Usage:

buld.gradle
```java
dependencies {
...
    implementation (name:'gomobilemodbus', ext:'aar')
}
```
import
```java
import gomobilemodbus.Gomobilemodbus;

```

### Examples Standard Modbus commands:
Read Coils
```java
String Modbusresult = Gomobilemodbus.modbusRequest("/dev/ttyS3", "46", "rc", "0", "1");
```

Read Holding Registers
```java
String Modbusresult = Gomobilemodbus.modbusRequest ("/dev/ttyS3", "46", "rh", "8", "");
```

Read Discrete Inputs
```java
String Modbusresult = Gomobilemodbus.modbusRequest("/dev/ttyS3", "46", "rd", "0", "1");
```

Write Single Coil
```java
String Modbusresult = Gomobilemodbus.modbusRequest("/dev/ttyS3", "46", "wc", "0", "1");
```

Write Single Holding Register
```java
String Modbusresult = Gomobilemodbus.modbusRequest ("/dev/ttyS3", "46", "wh", "9", "3");
```


### Proprietary commands for WB-MRM2-mini

Read Voltage
```java
String volt = Gomobilemodbus.modbusRequest("/dev/ttyS3", "46", "volt", "", "");
```

Read counter
```java
String Modbusresult = Gomobilemodbus.modbusRequest("/dev/ttyS3", "46", "count", "", "");
```

On Relay 1
```java
String modbusout = Gomobilemodbus.modbusRequest ("/dev/ttyS3", "46", "on", "0", "");
```
Off Relay 1
```java
String modbusout = Gomobilemodbus.modbusRequest ("/dev/ttyS3", "46", "off", "0", "");
```

On Relay 2
```java
String modbusout = Gomobilemodbus.modbusRequest ("/dev/ttyS3", "46", "on", "1", "");
```

Off Relay 2
```java
String modbusout = Gomobilemodbus.modbusRequest ("/dev/ttyS3", "46", "off", "1", "");
```
