package record

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

// Record records signal changes on the given pin and outputs the pulse and gap
// durations upon receiving SIGINT.
//
// Most remotes use the NEC protocol, which is documented at
// https://techdocs.altium.com/display/FPGA/NEC+Infrared+Transmission+Protocol.
func Record(pin uint8) {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	type Edge struct {
		V uint8
		T time.Time
	}
	xs := make([]Edge, 0, 2000)

	p := rpio.Pin(pin)
	p.Input()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	x := p.Read()
loop:
	for {
		select {
		case <-ch:
			rpio.Close()
			break loop
		default:
		}

		v := p.Read()
		if v == x {
			continue
		}
		x = v

		// Do as little as possible in this loop, or we'll miss
		// some edges. For instance, fmt.Println is way, way too
		// slow.
		xs = append(xs, Edge{
			T: time.Now(),
			V: uint8(v),
		})
	}

	if len(xs) < 2 {
		log.Println("No edges detected")
		return
	}

	for i := 1; i < len(xs); i++ {
		d := xs[i].T.Sub(xs[i-1].T)
		fmt.Printf("%d %6d\n", xs[i-1].V, int(d/time.Microsecond))
	}
}
