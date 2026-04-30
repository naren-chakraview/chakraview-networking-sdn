pub mod xdp;

pub struct EbpfConfig {
    pub interface: String,
    pub vni_map_size: usize,
    pub route_map_size: usize,
}

impl Default for EbpfConfig {
    fn default() -> Self {
        EbpfConfig {
            interface: "eth0".to_string(),
            vni_map_size: 100,
            route_map_size: 1000,
        }
    }
}

pub struct EbpfProgram {
    config: EbpfConfig,
}

impl EbpfProgram {
    pub fn new(config: EbpfConfig) -> Self {
        EbpfProgram { config }
    }

    pub fn load(&self) -> Result<(), String> {
        println!("Loading eBPF program for interface: {}", self.config.interface);
        println!("Route map size: {}", self.config.route_map_size);
        println!("VNI map size: {}", self.config.vni_map_size);

        Ok(())
    }

    pub fn get_stats(&self) -> Result<(u64, u64), String> {
        Ok((0, 0))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ebpf_config_defaults() {
        let config = EbpfConfig::default();
        assert_eq!(config.interface, "eth0");
        assert_eq!(config.route_map_size, 1000);
    }

    #[test]
    fn test_ebpf_program_creation() {
        let config = EbpfConfig::default();
        let program = EbpfProgram::new(config);
        assert!(program.load().is_ok());
    }
}
