# Send computer temperature to arduino
Send to arduino my setup temperature and display on lcd screen

## Dependencies
- libsensors-dev

## Start script when connect arduino
[blog](https://embarcados.com.br/utilizando-o-udev-para-criar-automacoes-com-porta-usb-no-linux/)
```
vim ~/.config/arduino.sh
```

Paste code on that new file
```
#!/bin/bash

echo "Arduino conectado!" > /tmp/arduino_connected.log
```

Edit udev rules
```
sudo vim /etc/udev/rules.d/99-arduino.rules
```

Paste the code
```
SUBSYSTEM=="tty", ATTRS{idVendor}=="2341", ATTRS{idProduct}=="0043", ACTION=="add", RUN+="<full-path-to-home>/.config/arduino.sh"
```

Reload udev rules
```
sudo udevadm control --reload-rules
```

When disconnect and reconnect the arcduino, runs
```
tail -f /tmp/arduino_connected.log
```

## Serial communication
[go serial](https://pkg.go.dev/go.bug.st/serial)

## My Setup
[k10temp-pci-00c3 (amd processor)](https://docs.kernel.org/hwmon/k10temp.html#description)