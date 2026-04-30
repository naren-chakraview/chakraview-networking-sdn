"""EVPN Route Types (RFC 7432)"""

from enum import Enum
from dataclasses import dataclass
from typing import Optional


class EVPNRouteType(Enum):
    ETHERNET_AD = 1
    MAC_IP = 2
    INCLUSIVE_MCAST = 3
    ETHERNET_SEGMENT = 4
    IP_PREFIX = 5


@dataclass
class EVPNRoute:
    route_type: EVPNRouteType
    route_distinguisher: str
    ethernet_tag: int
    esi: Optional[str] = None  # For Type 1, 4
    mac_address: Optional[str] = None  # For Type 2
    ip_address: Optional[str] = None  # For Type 2, 5
    next_hop: str = ""
    origin: str = "IGP"

    def __str__(self) -> str:
        return f"EVPN-{self.route_type.name}({self.route_distinguisher})"


def create_mac_ip_route(rd: str, eth_tag: int, mac: str,
                        ip: str, nh: str) -> EVPNRoute:
    """Create Type 2 (MAC/IP) route"""
    return EVPNRoute(
        route_type=EVPNRouteType.MAC_IP,
        route_distinguisher=rd,
        ethernet_tag=eth_tag,
        mac_address=mac,
        ip_address=ip,
        next_hop=nh
    )


def create_ip_prefix_route(rd: str, eth_tag: int,
                           prefix: str, nh: str) -> EVPNRoute:
    """Create Type 5 (IP Prefix) route"""
    return EVPNRoute(
        route_type=EVPNRouteType.IP_PREFIX,
        route_distinguisher=rd,
        ethernet_tag=eth_tag,
        ip_address=prefix,
        next_hop=nh
    )
