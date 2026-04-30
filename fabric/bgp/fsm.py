"""BGP Finite State Machine (RFC 4271)"""

from enum import Enum
from typing import Optional, List
from dataclasses import dataclass
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class BGPState(Enum):
    IDLE = "IDLE"
    CONNECT = "CONNECT"
    ACTIVE = "ACTIVE"
    OPENSENT = "OPENSENT"
    OPENCONFIRM = "OPENCONFIRM"
    ESTABLISHED = "ESTABLISHED"


class BGPEvent(Enum):
    START = "Start"
    STOP = "Stop"
    TRANSPORT_CONN_OPEN = "TransportConnOpen"
    TRANSPORT_CONN_FAIL = "TransportConnFail"
    TRANSPORT_CONN_CLOSE = "TransportConnClose"
    BGP_OPEN = "BGPOpen"
    BGP_HEADER_ERR = "BGPHeaderErr"
    BGP_OPEN_MSG_ERR = "BGPOpenMsgErr"
    KEEPALIVE_MSG = "KeepaliveMsg"
    UPDATE_MSG = "UpdateMsg"


@dataclass
class BGPPeer:
    asn: int
    router_id: str
    address: str
    state: BGPState = BGPState.IDLE
    hold_time: int = 180
    keepalive_interval: int = 60


class BGPStateMachine:
    def __init__(self, local_asn: int, router_id: str):
        self.local_asn = local_asn
        self.router_id = router_id
        self.state = BGPState.IDLE
        self.peers: dict = {}

    def add_peer(self, peer_addr: str, peer_asn: int, peer_router_id: str):
        self.peers[peer_addr] = BGPPeer(
            asn=peer_asn,
            router_id=peer_router_id,
            address=peer_addr
        )
        logger.info(f"BGP peer added: {peer_addr} (ASN {peer_asn})")

    def process_event(self, event: BGPEvent, peer_addr: Optional[str] = None) -> bool:
        if peer_addr and peer_addr in self.peers:
            peer = self.peers[peer_addr]
            return self._transition(peer, event)
        return False

    def _transition(self, peer: BGPPeer, event: BGPEvent) -> bool:
        old_state = peer.state

        if peer.state == BGPState.IDLE and event == BGPEvent.START:
            peer.state = BGPState.CONNECT
        elif peer.state == BGPState.CONNECT and event == BGPEvent.TRANSPORT_CONN_OPEN:
            peer.state = BGPState.OPENSENT
        elif peer.state == BGPState.OPENSENT and event == BGPEvent.BGP_OPEN:
            peer.state = BGPState.OPENCONFIRM
        elif peer.state == BGPState.OPENCONFIRM and event == BGPEvent.KEEPALIVE_MSG:
            peer.state = BGPState.ESTABLISHED
        elif event == BGPEvent.STOP or event == BGPEvent.TRANSPORT_CONN_CLOSE:
            peer.state = BGPState.IDLE

        if old_state != peer.state:
            logger.info(f"BGP peer {peer.address}: {old_state.value} -> {peer.state.value}")
            return True
        return False

    def get_peer_state(self, peer_addr: str) -> Optional[BGPState]:
        if peer_addr in self.peers:
            return self.peers[peer_addr].state
        return None

    def list_established_peers(self) -> List[str]:
        return [addr for addr, peer in self.peers.items()
                if peer.state == BGPState.ESTABLISHED]
