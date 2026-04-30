.PHONY: build test clean

build:
	cd dpdk && make
	cd ebpf && cargo build --release

test:
	cd dpdk && make test
	cd ebpf && cargo test

clean:
	cd dpdk && make clean
	cd ebpf && cargo clean
