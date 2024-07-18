package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/md14454/gosensors"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func main() {
	// serialCommunication()
	// hardwareSensors()
	nvidiaSensors()
}

func nvidiaSensors() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(ret))
	}

	for i := 0; i < count; i++ {
		device, ret := nvml.DeviceGetHandleByIndex(i)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, nvml.ErrorString(ret))
		}

		uuid, ret := device.GetUUID()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get uuid of device at index %d: %v", i, nvml.ErrorString(ret))
		}

		temp, _ := device.GetTemperature(nvml.TEMPERATURE_GPU)
		fmt.Printf("%v\t%+v\n", uuid, temp)
	}
}

func hardwareSensors() {
	gosensors.Init()
	defer gosensors.Cleanup()

	chips := gosensors.GetDetectedChips()

	for i := 0; i < len(chips); i++ {
		chip := chips[i]

		fmt.Printf("%v\n", chip)
		fmt.Printf("Adapter: %v\n", chip.AdapterName())

		features := chip.GetFeatures()

		for j := 0; j < len(features); j++ {
			feature := features[j]

			fmt.Printf("%v ('%v'): %.1f\n", feature.Name, feature.GetLabel(), feature.GetValue())

			subfeatures := feature.GetSubFeatures()

			for k := 0; k < len(subfeatures); k++ {
				subfeature := subfeatures[k]

				fmt.Printf("  %v: %.1f\n", subfeature.Name, subfeature.GetValue())
			}
		}

		fmt.Printf("\n")
	}
}

func serialCommunication() {
	// read available ports
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %+v\n", port)
		// metadata
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
		}
	}

	sensor := exec.Command("sensors")
	resp, _ := sensor.Output()
	fmt.Printf("Sensors response:\n%+v", string(resp))

	// open communication
	device := ports[0].Name
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found port: %+v\n", port)

	// IO communication
	n, err := port.Write([]byte("10,20,30\n\r"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))
	}
}
