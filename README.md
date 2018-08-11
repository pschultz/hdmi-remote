HTTP server for controlling my HDMI switch over IR. Runs on a Raspberry PI.

# Setup

    user@localhost$ GOARCH=arm go build
    user@localhost$ scp hdmi-remote root@pi:/usr/local/bin/hdmi-remote
    user@localhost$ scp hdmi-remote.service root@pi:/etc/systemd/system
    user@localhost$ scp -r _irsling root@pi:

    root@pi# cd _irsling
    root@pi# ( cd vendor/github.com/joan2937/pigpio; make && make install )
    root@pi# make install
    root@pi# systemctl enable hdmi-remote
    root@pi# systemctl start hdmi-remote

Special thanks to @bschwind. His [blog post about IR][bschwind] and [irsling-er
library][irslinger] were an incredible help and enabled me to finish this project
in less than a day.

[bschwind]: https://blog.bschwind.com/2016/05/29/sending-infrared-commands-from-a-raspberry-pi-without-lirc/
[irslinger]: https://github.com/bschwind/ir-slinger
