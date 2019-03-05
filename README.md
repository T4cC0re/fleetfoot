```
  __ _           _    __            _
 / _| |         | |  / _|          | |
| |_| | ___  ___| |_| |_ ___   ___ | |_
|  _| |/ _ \/ _ \ __|  _/ _ \ / _ \| __|
| | | |  __/  __/ |_| || (_) | (_) | |_
|_| |_|\___|\___|\__|_| \___/ \___/ \__|

   ...Arch Linux based routing made easy
```
*This repo is primarily hosted on [gitlab.com](https://gitlab.com/T4cC0re/sonicrainboom). Please go there to raise issues or contribute.*


#### Why?

...because why the heck not? No, in all seriousness. I was looking for a more
powerful router that could cope with PPPoE, VLANs, lots of Gbit/s in a mixed
copper, fiber and InfiniBand network and massive expandability, because I like
to hack around with stuff...

Needless to say, I did not really find something that flipped me on.

At first I looked at VyOS/Vyatta (the same stuff Ubiquitis Edge OS is built
ontop of), but I soon realized, that it was not as extensible as I hoped.

Introducing __*fleetfoot*__!

I really like Arch Linux for its very up2date system and ecosystem, so I went
that route, but instead of building ontop of netctl I choose to adapt
systemd-networkd so that this might be ported to other OSs should someone want
to.

Also everybody wants to say "I run an Arch Linux server!" :) Here you can
over-deliver and say "I run an Arch Linux router!". Even better IMO

#### Features

- 20% cooler than all other routers
- Boots in 10 seconds flat! \*
  - \* on an Intel Core i5-6500T + UEFI fast-boot on a ASRock B150M Pro4V
  - Powered by systemd-boot
- Routing (duh!)
- DHCP
  - Static leases
  - Range per bridge
- DNS
  - To come: DNS-over-TLS/DNS-over-HTTPS via local proxy behind dnsmasq
  - DNSSEC foo
- fiber!
- InfiniBand (soon!)
- PPPoE
  - With hooks!
  - Port forwarding with hairpinning
  - Drop unsolicited packets by default
  - mss-clamping, so your clients will not be screwed by the path-MTU
- VLANs - Everything is a bridge
  - Except PPPoE and maybe VPNs...
  - Attach any interface to any VLAN or multiple at once!
- Expandability up the bazoo
  - It's based on Arch Linux after all

#### How to try it?

Bad news... I do not have a proper guide to install it right now.
The easiest way is, to install Arch, python and openssh onto a box, configure
your `inventory.yml` and let Ansible to the rest.

Speaking of which, I also do not have a proper sample `inventory.yml`, yet.
This is because I will re-do the Ansible stuff as a tiny golang daemon doing
all the dirty work, an also listen on the hooks to make it a much more
streamlined experience.

This (or more like the initial parts of it) can be found in the `fleetfootd`
directory.
