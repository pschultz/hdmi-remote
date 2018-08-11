#include <stdio.h>
#include "irslinger.h"

// compile with `gcc -o irsling irsling.c -Ipath/to/ir-slinger -lm -lpigpio -pthread -lrt` 

int main(int argc, char *argv[])
{
	if (argc != 2) {
		fprintf(stderr, "Usage: %s <BINARY CODE>\n", argv[0]);
		fprintf(stderr, "\n");
		fprintf(stderr, "BINARY CODE is the binary representation of the data to send, typically\n");
		fprintf(stderr, "32 bits for NEC (addr || ~addr || command || ~command).\n");
		fprintf(stderr, "\n");
		fprintf(stderr, "Example: %s 00000000111111111010101001010101", argv[0]);

		return -1;
	}

	uint32_t outPin = 18;            // The Broadcom pin number the signal will be sent on
	int frequency = 38000;           // The frequency of the IR signal in Hz
	double dutyCycle = 0.5;          // The duty cycle of the IR signal. 0.5 means for every cycle,
									 // the LED will turn on for half the cycle
									 // time, and off the other half
	int leadingPulseDuration = 9000; // The duration of the beginning pulse in microseconds
	int leadingGapDuration = 4500;   // The duration of the gap in microseconds after the leading pulse
	int onePulse = 562;              // The duration of a pulse in microseconds when sending a logical 1
	int zeroPulse = 562;             // The duration of a pulse in microseconds when sending a logical 0
	int oneGap = 1688;               // The duration of the gap in microseconds when sending a logical 1
	int zeroGap = 562;               // The duration of the gap in microseconds when sending a logical 0
	int sendTrailingPulse = 1;       // 1 = Send a trailing pulse with duration equal to "onePulse"
	                                 // 0 = Don't send a trailing pulse

	return irSling(
		outPin,
		frequency,
		dutyCycle,
		leadingPulseDuration,
		leadingGapDuration,
		onePulse,
		zeroPulse,
		oneGap,
		zeroGap,
		sendTrailingPulse,
		argv[1]
	);
}
