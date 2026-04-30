"""VXLAN Tunnel Encapsulation and Management"""

from dataclasses import dataclass
from typing import Dict, Optional
import logging

logger = logging.getLogger(__name__)


@dataclass
class VXLANTunnel:
    tunnel_id: str
    local_ip: str
    remote_ip: str
    vni: int
    mtu: int = 1500
    active: bool = True

    def encapsulate(self, inner_packet: bytes) -> bytes:
        """Encapsulate packet in VXLAN header"""
        vxlan_header = self._build_vxlan_header()
        outer_ip = self._build_outer_ip()
        outer_udp = self._build_outer_udp()
        outer_eth = self._build_outer_eth()

        return outer_eth + outer_ip + outer_udp + vxlan_header + inner_packet

    def _build_vxlan_header(self) -> bytes:
        """Build 8-byte VXLAN header"""
        flags = 0x08  # I bit set
        reserved = 0x000000
        vni_bytes = self.vni.to_bytes(3, 'big')
        return bytes([flags]) + reserved.to_bytes(3, 'big') + vni_bytes + b'\x00'

    def _build_outer_eth(self) -> bytes:
        """Simplified outer Ethernet header"""
        return b'\xff' * 6 + b'\x00' * 6 + b'\x08\x00'

    def _build_outer_ip(self) -> bytes:
        """Simplified outer IP header"""
        version_ihl = 0x45
        dscp_ecn = 0x00
        total_length = 20 + 8 + 8 + 100  # Placeholder
        return (bytes([version_ihl, dscp_ecn]) +
                total_length.to_bytes(2, 'big') +
                b'\x00' * 4 +  # ID, flags, frag offset
                bytes([64, 17]) +  # TTL=64, Protocol=UDP
                b'\x00' * 2 +  # Checksum (optional)
                bytes(map(int, self.local_ip.split('.'))) +
                bytes(map(int, self.remote_ip.split('.'))))

    def _build_outer_udp(self) -> bytes:
        """Simplified outer UDP header"""
        src_port = 4789
        dst_port = 4789
        length = 8 + 8 + 100  # Placeholder
        return (src_port.to_bytes(2, 'big') +
                dst_port.to_bytes(2, 'big') +
                length.to_bytes(2, 'big') +
                b'\x00\x00')  # Checksum optional


class TunnelManager:
    def __init__(self):
        self.tunnels: Dict[str, VXLANTunnel] = {}

    def create_tunnel(self, tunnel_id: str, local_ip: str,
                      remote_ip: str, vni: int) -> VXLANTunnel:
        tunnel = VXLANTunnel(
            tunnel_id=tunnel_id,
            local_ip=local_ip,
            remote_ip=remote_ip,
            vni=vni
        )
        self.tunnels[tunnel_id] = tunnel
        logger.info(f"VXLAN tunnel created: {tunnel_id} (VNI {vni})")
        return tunnel

    def get_tunnel(self, tunnel_id: str) -> Optional[VXLANTunnel]:
        return self.tunnels.get(tunnel_id)

    def delete_tunnel(self, tunnel_id: str) -> bool:
        if tunnel_id in self.tunnels:
            del self.tunnels[tunnel_id]
            logger.info(f"VXLAN tunnel deleted: {tunnel_id}")
            return True
        return False

    def list_tunnels(self) -> list:
        return list(self.tunnels.values())
