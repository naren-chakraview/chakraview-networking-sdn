"""EVPN Route Management"""

from typing import Dict, List, Optional
from .types import EVPNRoute, EVPNRouteType
import logging

logger = logging.getLogger(__name__)


class EVPNRouteManager:
    def __init__(self):
        self.routes: Dict[str, List[EVPNRoute]] = {}
        self.rib: Dict[str, EVPNRoute] = {}

    def announce_route(self, route: EVPNRoute) -> None:
        """Announce new EVPN route"""
        key = f"{route.route_type.name}:{route.route_distinguisher}"

        if key not in self.routes:
            self.routes[key] = []

        self.routes[key].append(route)
        self.rib[key] = route

        logger.info(f"EVPN route announced: {route}")

    def withdraw_route(self, route_type: EVPNRouteType, rd: str) -> bool:
        """Withdraw EVPN route"""
        key = f"{route_type.name}:{rd}"

        if key in self.rib:
            del self.rib[key]
            if key in self.routes:
                del self.routes[key]
            logger.info(f"EVPN route withdrawn: {key}")
            return True
        return False

    def get_mac_ip_routes(self) -> List[EVPNRoute]:
        """Get all Type 2 (MAC/IP) routes"""
        return [r for routes in self.routes.values() for r in routes
                if r.route_type == EVPNRouteType.MAC_IP]

    def get_ip_prefix_routes(self) -> List[EVPNRoute]:
        """Get all Type 5 (IP Prefix) routes"""
        return [r for routes in self.routes.values() for r in routes
                if r.route_type == EVPNRouteType.IP_PREFIX]

    def get_rib(self) -> Dict[str, EVPNRoute]:
        """Get RIB (Routing Information Base)"""
        return dict(self.rib)

    def lookup_route(self, rd: str) -> Optional[EVPNRoute]:
        """Lookup route by RD"""
        for key, route in self.rib.items():
            if rd in key:
                return route
        return None
