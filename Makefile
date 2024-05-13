# cfs-gui Makefile

phony := build 
phony += all
phony += backend
phony += frontend
phony += clean

default: build

build: build.sh
	bash build.sh --all

all:
	bash build.sh --all

backend:
	bash build.sh --backend

frontend:
	bash build.sh --frontend

clean: bin
	bash -rf bin

.PHONY: $(phony)
