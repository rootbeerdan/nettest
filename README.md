# nettest
Nettest is a small utility that will run ICMP echos to specified endpoints over IPv4 and IPv6, and expose ICMP metrics to Prometheus.

## How to use

Download the binary, and run it as a background service. Metrics will be availible over localhost:8081/metrics

## Compatibility

Tested on Windows 11, however it should work on Linux and macOS as a privileged application, as ICMP is considered low level access on UNIX systems.

## Issues & feature requests

Feel free to open up a ticket, but this was a one-off thing and I'm pretty busy these days.

## Known issues

- Program will show endpoint success if either IPv4 or IPv6 endpoint is working, which may not be desirable. Remediated by using single stack endpoints like ipv4.google.com and ipv6.google.com.
