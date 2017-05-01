package main

import (
	"github.com/gitaiqaq/serial"
	"io/ioutil"
	"time"
	"github.com/elastic/beats/libbeat/common"
	"strings"
	"strconv"
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

func is_Frame(b []byte) bool {
	if (len(b) < 54) {
		return false
	}
	return is_Number(b[0]) && b[1] == 124
}

type Frame struct {
	version 		int
	frameType 		int
	frameSubType 	int
	time 		 	common.Time
	chipId 			string
	rssi			string
	channel			int
	receiverMAC		string
	senderMAC		string
	ssid			string
}

func main() {
        c := &serial.Config{Name: findArduino(), Baud: 115200}
        s, err := serial.OpenPort(c)
        if err != nil {
                log.Fatal(err)
        }

		scanner := bufio.NewScanner(s.File())
	    for scanner.Scan() {
	    	if(is_Frame(scanner.Bytes())){
		    	tokens := strings.Split(scanner.Text(), "|")
		    	version, err	:= strconv.Atoi(tokens[0])
		    	if err != nil {
		    		continue
		    	}
		    	frameType, err	:= strconv.Atoi(tokens[1])
		    	if err != nil {
		    		continue
		    	}
		    	frameSubType, err	:= strconv.Atoi(tokens[2])
		    	if err != nil {
		    		continue
		    	}
		    	channel, err	:= strconv.Atoi(tokens[6])
		    	if err != nil {
		    		channel = 0
		    	}
		    	frame := Frame{
		    		version:		version,
					frameType:		frameType,
					frameSubType:	frameSubType,
					time:			common.Time(time.Now()),
					chipId:			tokens[3],
					rssi:			tokens[5],
					channel:		channel,
					receiverMAC:	tokens[7],
					senderMAC:		tokens[8],
					ssid:			tokens[9],
		    	}
		    	fmt.Println(frame)
	    	}
	    }

	    if err := scanner.Err(); err != nil {
	        log.Fatal(err)
	    }
}
