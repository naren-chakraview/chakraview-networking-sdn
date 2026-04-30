"""BGP Route Table and Advertisement"""

from dataclasses import dataclass
from typing import Optional, List, Dict
import logging

logger = logging.getLogger(__name__)


@dataclass
class BGPRoute:
    destination: str
    next_hop: str
    as_path: str
    origin: str = "IGP"
    local_pref: int = 100
    med: int = 0


class RouteTable:
    def __init__(self):
        self.routes: Dict[str, BGPRoute] = {}
        self.advertised: Dict[str, List[str]] = {}

    def add_route(self, destination: str, next_hop: str, as_path: str) -> None:
        route = BGPRoute(
            destination=destination,
            next_hop=next_hop,
            as_path=as_path
        )
        self.routes[destination] = route
        logger.info(f"Route learned: {destination} via {next_hop}")

    def advertise_route(self, destination: str, to_peer: str) -> bool:
        if destination not in self.routes:
            logger.warning(f"Cannot advertise unknown route: {destination}")
            return False

        if destination not in self.advertised:
            self.advertised[destination] = []

        if to_peer not in self.advertised[destination]:
            self.advertised[destination].append(to_peer)
            logger.info(f"Route advertised: {destination} to {to_peer}")
            return True
        return False

    def withdraw_route(self, destination: str) -> None:
        if destination in self.routes:
            del self.routes[destination]
        if destination in self.advertised:
            del self.advertised[destination]
        logger.info(f"Route withdrawn: {destination}")

    def get_route(self, destination: str) -> Optional[BGPRoute]:
        return self.routes.get(destination)

    def list_routes(self) -> List[BGPRoute]:
        return list(self.routes.values())

    def get_best_route(self, destination: str) -> Optional[BGPRoute]:
        return self.routes.get(destination)
