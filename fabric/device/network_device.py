"""Simulated SDN-enabled Network Device"""

from enum import Enum
from dataclasses import dataclass
from typing import Dict, List, Optional
import logging

logger = logging.getLogger(__name__)


class DeviceRole(Enum):
    LEAF = "leaf"
    SPINE = "spine"
    BORDER = "border"


@dataclass
class Interface:
    name: str
    ip_address: str
    vlan: int = 0
    mtu: int = 1500
    enabled: bool = True


class NetworkDevice:
    def __init__(self, device_id: str, asn: int, router_id: str,
                 role: DeviceRole, mgmt_ip: str):
        self.device_id = device_id
        self.asn = asn
        self.router_id = router_id
        self.role = role
        self.mgmt_ip = mgmt_ip

        self.interfaces: Dict[str, Interface] = {}
        self.routing_table: Dict[str, str] = {}
        self.vxlan_tunnels: Dict[str, dict] = {}
        self.bgp_peers: List[str] = []
        self.mac_table: Dict[str, str] = {}

        self.packets_forwarded = 0
        self.packets_dropped = 0

    def add_interface(self, name: str, ip: str, vlan: int = 0) -> Interface:
        """Add interface to device"""
        iface = Interface(name=name, ip_address=ip, vlan=vlan)
        self.interfaces[name] = iface
        logger.info(f"{self.device_id}: Interface {name} added ({ip})")
        return iface

    def add_route(self, destination: str, next_hop: str) -> None:
        """Add route to routing table"""
        self.routing_table[destination] = next_hop
        logger.info(f"{self.device_id}: Route added {destination} -> {next_hop}")

    def forward_packet(self, dest_ip: str) -> bool:
        """Simulate packet forwarding"""
        if dest_ip in self.routing_table:
            self.packets_forwarded += 1
            return True
        else:
            self.packets_dropped += 1
            return False

    def add_bgp_peer(self, peer_addr: str) -> None:
        """Add BGP peer"""
        self.bgp_peers.append(peer_addr)
        logger.info(f"{self.device_id}: BGP peer added {peer_addr}")

    def create_vxlan_tunnel(self, tunnel_id: str, remote_ip: str, vni: int) -> None:
        """Create VXLAN tunnel"""
        self.vxlan_tunnels[tunnel_id] = {
            'remote_ip': remote_ip,
            'vni': vni,
            'active': True
        }
        logger.info(f"{self.device_id}: VXLAN tunnel created {tunnel_id} to {remote_ip}")

    def learn_mac(self, mac: str, interface: str) -> None:
        """Learn MAC address on interface"""
        self.mac_table[mac] = interface
        logger.info(f"{self.device_id}: MAC learned {mac} on {interface}")

    def get_stats(self) -> Dict:
        """Get device statistics"""
        return {
            'device_id': self.device_id,
            'role': self.role.value,
            'packets_forwarded': self.packets_forwarded,
            'packets_dropped': self.packets_dropped,
            'routes': len(self.routing_table),
            'tunnels': len(self.vxlan_tunnels),
            'mac_table_size': len(self.mac_table),
            'bgp_peers': len(self.bgp_peers)
        }
