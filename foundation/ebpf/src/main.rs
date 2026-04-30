use networking_sdn_ebpf::{EbpfConfig, EbpfProgram};

fn main() {
    println!("eBPF Program Loader");

    let config = EbpfConfig::default();
    let program = EbpfProgram::new(config);

    match program.load() {
        Ok(()) => {
            println!("eBPF program loaded successfully");

            match program.get_stats() {
                Ok((fwd, drop)) => {
                    println!("Initial stats: forwarded={}, dropped={}", fwd, drop);
                }
                Err(e) => {
                    eprintln!("Failed to get stats: {}", e);
                }
            }
        }
        Err(e) => {
            eprintln!("Failed to load eBPF program: {}", e);
            std::process::exit(1);
        }
    }
}
