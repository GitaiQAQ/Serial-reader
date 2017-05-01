package main

import (
	"github.com/gitaiqaq/serial"
	"io/ioutil"
	"strings"
	"bufio"
	"fmt"
	"log"
)

// findArduino looks for the file that represents the Arduino
// serial connection. Returns the fully qualified path to the
// device if we are able to find a likely candidate for an
// Arduino, otherwise an empty string if unable to find
// something that 'looks' like an Arduino device.
func findArduino() string {
	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyUSB") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}

func is_Number(b byte) bool {
	return b >= 48 && b <= 57
}

func is_Line(b []byte) bool {
	return is_Number(b[0]) && b[1] == 124
}

func main() {
        c := &serial.Config{Name: findArduino(), Baud: 115200}
        s, err := serial.OpenPort(c)
        if err != nil {
                log.Fatal(err)
        }

		scanner := bufio.NewScanner(s.File())
	    for scanner.Scan() {
	    	if (is_Line(scanner.Bytes())) {
	    		fmt.Print(scanner.Bytes())
	    	}else{
	    		//fmt.Println(scanner.Text())
	    	}
	    }

	    if err := scanner.Err(); err != nil {
	        log.Fatal(err)
	    }
}
