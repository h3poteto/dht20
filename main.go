package main

import (
	"log"
	"math"
	"time"

	"github.com/d2r2/go-i2c"
)

// Refs: https://craft-gogo.com/raspberry-pi-dht20/
// Refs: https://s-design-tokyo.com/use-dht20-raspberrypi/
// Refs: https://hatakekara.com/dht20-arduino/

func main() {
	// Create new connection to I2C bus on 2 line with address 0x38
	i2c, err := i2c.NewI2C(0x38, 1)
	if err != nil {
		log.Fatal(err)
	}
	// Free I2C connection on exit
	defer i2c.Close()

	// Need to wait 100ms after launched
	time.Sleep(100 * time.Millisecond)
	var initial byte = 0x71
	ret, err := i2c.ReadRegU8(initial)
	if err != nil {
		log.Fatal(err)
	}
	// ret should be 28(0x1c)
	log.Printf("init code %d", ret)

	time.Sleep(10 * time.Millisecond)
	// Start measure
	_, err = i2c.WriteBytes([]byte{0x00, 0xAC, 0x33, 0x00})
	if err != nil {
		log.Fatal(err)
	}
	// Need to wait after sending ac3300
	time.Sleep(80 * time.Millisecond)
	dat := make([]byte, 7)
	_, err = i2c.ReadBytes(dat)
	if err != nil {
		log.Fatal(err)
	}
	// byte is uint8 (8bits), it is not enough to shift
	// So cast []uint8 to []uint32
	var long_dat []uint32
	for _, d := range dat {
		long_dat = append(long_dat, uint32(d))
	}

	// Get humidity and tempreature data
	hum := long_dat[1]<<12 | long_dat[2]<<4 | ((long_dat[3] & 0xF0) >> 4)
	tmp := ((long_dat[3] & 0x0F) << 16) | long_dat[4]<<8 | long_dat[5]

	// Calcurate real data
	real_hum := float64(hum) / math.Pow(2, 20) * 100
	real_tmp := float64(tmp)/math.Pow(2, 20)*200 - 50

	log.Printf("hum: %f", real_hum)
	log.Printf("tmp: %f", real_tmp)
}
