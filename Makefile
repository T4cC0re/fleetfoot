.PHONY: clean all

all: bin/miniupnpd bin/fleetfootd bin/fleetfootctl

bin:
	mkdir -p bin

bin/miniupnpd: bin miniupnp/miniupnpd/miniupnpd
	cp -rvp miniupnp/miniupnpd/miniupnpd bin/

bin/fleetfootd: bin
	go build -o bin/fleetfootd ./daemon

bin/fleetfootctl: bin
	go build -o bin/fleetfootctl ./ctl

miniupnp:
	git clone -c advice.detachedHead=false --depth 1 -b miniupnpd_2_1 https://framagit.org/miniupnp/miniupnp.git

miniupnp/miniupnpd/miniupnpd: miniupnp
	cd miniupnp/miniupnpd; make -j3 -f Makefile.linux

clean:
	rm -rf miniupnp bin
