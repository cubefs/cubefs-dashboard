# cfs-gui Makefile

phony := build 
phony += clean

build: build.sh
	sh build.sh

clean: bin
	rm -rf bin

.PHONY: $(phony)
