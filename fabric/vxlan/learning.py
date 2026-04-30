"""Dynamic MAC Learning for VXLAN"""

from dataclasses import dataclass
from typing import Dict, Optional, Tuple
import logging

logger = logging.getLogger(__name__)


@dataclass
class MACEntry:
    mac_address: str
    tunnel_id: str
    learned_from_vni: int
    age: int = 0


class MACLearningTable:
    def __init__(self):
        self.mac_table: Dict[str, MACEntry] = {}
        self.max_age = 300  # 5 minutes

    def learn_mac(self, mac_address: str, tunnel_id: str, vni: int) -> None:
        """Learn MAC address via tunnel and VNI"""
        entry = MACEntry(
            mac_address=mac_address,
            tunnel_id=tunnel_id,
            learned_from_vni=vni
        )
        self.mac_table[mac_address] = entry
        logger.info(f"MAC learned: {mac_address} via {tunnel_id} (VNI {vni})")

    def lookup_mac(self, mac_address: str) -> Optional[Tuple[str, int]]:
        """Lookup MAC address, return (tunnel_id, vni)"""
        entry = self.mac_table.get(mac_address)
        if entry:
            return (entry.tunnel_id, entry.learned_from_vni)
        return None

    def age_out_macs(self) -> int:
        """Age out stale entries, return count removed"""
        expired = [mac for mac, entry in self.mac_table.items()
                   if entry.age >= self.max_age]
        for mac in expired:
            del self.mac_table[mac]
            logger.info(f"MAC aged out: {mac}")
        return len(expired)

    def flush_tunnel(self, tunnel_id: str) -> int:
        """Remove all MACs learned via tunnel"""
        to_remove = [mac for mac, entry in self.mac_table.items()
                     if entry.tunnel_id == tunnel_id]
        for mac in to_remove:
            del self.mac_table[mac]
        logger.info(f"Flushed {len(to_remove)} MACs from {tunnel_id}")
        return len(to_remove)

    def get_table(self) -> Dict[str, MACEntry]:
        return dict(self.mac_table)
