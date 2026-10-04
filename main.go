package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/warthog618/go-gpiocdev"
)

const FanControl = 18
const FanOff = 1
const FanOn = 0

type Driver struct {
	line *gpiocdev.Line

	counterLowTemp  int
	counterHighTemp int
	counterLowCpu   int
	counterHighCpu  int
	counterStatus   int
	fanStatus       bool
}

func NewDriver() (*Driver, error) {
	line, err := gpiocdev.RequestLine(
		"gpiochip0",
		FanControl,
		gpiocdev.AsOutput(1),
		gpiocdev.AsActiveLow,
		gpiocdev.WithConsumer("fanshim-bare"),
	)
	if err != nil {
		return nil, fmt.Errorf("request fan GPIO line %d: %w", FanControl, err)
	}

	return &Driver{
		line: line,
	}, nil
}

func (drv *Driver) Close() {
	drv.line.Close()
}

func (drv *Driver) SetFan(val int) error {
	if val == FanOn {
		drv.fanStatus = true
	} else {
		drv.fanStatus = false
	}

	if err := drv.line.SetValue(val); err != nil {
		return fmt.Errorf("couldn't set fan to %d: %w", val, err)
	}
	return nil
}

func (drv *Driver) getSensorTemp(sensorKey string) (float64, error) {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return 0, nil
	}

	for _, temp := range temps {
		if temp.SensorKey == sensorKey {
			return temp.Temperature, nil
		}
	}

	return 0, fmt.Errorf("sensorKey \"%s\" not found", sensorKey)
}

func (drv *Driver) RunCheck() error {
	cpuTemp, err := drv.getSensorTemp("cpu_thermal")
	if err != nil {
		return err
	}

	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	cpuPercent := cpuPercents[0]

	if cpuTemp >= 65 {
		drv.counterHighTemp++
		drv.counterLowTemp = 0
	} else if cpuTemp < 55 {
		drv.counterHighTemp = 0
		drv.counterLowTemp++
	}
	if cpuPercent > 70 {
		drv.counterHighCpu++
		drv.counterLowCpu = 0
	} else if cpuPercent < 50 {
		drv.counterHighCpu = 0
		drv.counterLowCpu++
	}

	if (drv.counterHighTemp > 0 || drv.counterHighCpu > 10) && !drv.fanStatus {
		fmt.Printf("fan on: %.1f°C, %3.2f%%\n", cpuTemp, cpuPercent)
		drv.SetFan(FanOn)
	} else if drv.counterLowTemp > 5 && drv.counterLowCpu > 5 && drv.fanStatus {
		fmt.Printf("fan off: %.1f°C, %3.2f%%\n", cpuTemp, cpuPercent)
		drv.SetFan(FanOff)
	}

	return nil
}

func main() {
	drv, err := NewDriver()
	if err != nil {
		panic(err)
	}
	defer drv.Close()
	drv.SetFan(FanOn)

	for {
		if err := drv.RunCheck(); err != nil {
			drv.SetFan(FanOn)
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}
