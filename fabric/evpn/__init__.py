"""EVPN (Ethernet VPN) route types and handling"""

from .types import EVPNRouteType, EVPNRoute
from .routes import EVPNRouteManager

__all__ = ['EVPNRouteType', 'EVPNRoute', 'EVPNRouteManager']
