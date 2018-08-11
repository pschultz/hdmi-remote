irsling: irsling.c vendor/github.com/bschwind/ir-slinger/irslinger.h
	gcc -o $@ -I vendor/github.com/bschwind/ir-slinger -lm -lpigpio -pthread -lrt irsling.c

install: irsling
	cp irsling /usr/local/bin/irsling
