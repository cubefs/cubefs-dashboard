# cfs-gui Makefile

phony := build 
phony += all
phony += backend
phony += frontend
phony += clean

default: build

build: build.sh
	sh build.sh --all

all:
	sh build.sh --all

backend:
	sh build.sh --backend

frontend:
	sh build.sh --frontend

clean: bin
	rm -rf bin

.PHONY: $(phony)
