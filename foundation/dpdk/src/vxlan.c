#include "vxlan.h"
#include <string.h>
#include <arpa/inet.h>

void vxlan_init(vxlan_state_t *state) {
    memset(state, 0, sizeof(vxlan_state_t));
}

bool vxlan_add_tunnel(vxlan_state_t *state, uint32_t tunnel_id, uint32_t local_ip,
                      uint32_t remote_ip, uint32_t vni) {
    if (state->tunnel_count >= MAX_VXLAN_TUNNELS) {
        return false;
    }

    vxlan_tunnel_t *tunnel = &state->tunnels[state->tunnel_count];
    tunnel->tunnel_id = tunnel_id;
    tunnel->local_ip = local_ip;
    tunnel->remote_ip = remote_ip;
    tunnel->vni = vni;
    tunnel->active = true;

    state->tunnel_count++;
    return true;
}

static vxlan_tunnel_t *find_tunnel(vxlan_state_t *state, uint32_t tunnel_id) {
    for (int i = 0; i < state->tunnel_count; i++) {
        if (state->tunnels[i].tunnel_id == tunnel_id && state->tunnels[i].active) {
            return &state->tunnels[i];
        }
    }
    return NULL;
}

/* Build outer UDP/VXLAN headers (simplified; real DPDK would use mbuf) */
bool vxlan_encapsulate(vxlan_state_t *state, uint32_t tunnel_id,
                       const uint8_t *input_pkt, uint32_t input_len,
                       uint8_t *output_pkt, uint32_t *output_len) {
    vxlan_tunnel_t *tunnel = find_tunnel(state, tunnel_id);
    if (!tunnel) {
        return false;
    }

    /* Header sizes: Eth(14) + IP(20) + UDP(8) + VXLAN(8) */
    uint32_t header_size = 14 + 20 + 8 + 8;
    if (header_size + input_len > 65535) {
        return false; /* Packet too large */
    }

    /* Copy inner packet */
    memcpy(output_pkt + header_size, input_pkt, input_len);

    /* Build VXLAN header (UDP payload) */
    vxlan_hdr_t *vxlan = (vxlan_hdr_t *)(output_pkt + header_size - 8);
    vxlan->flags = 0x08; /* I bit set (VNI valid) */
    vxlan->reserved1[0] = vxlan->reserved1[1] = vxlan->reserved1[2] = 0;
    vxlan->vni = htonl((tunnel->vni << 8) & 0xFFFFFF00);

    /* Outer IP header (simplified) */
    uint8_t *ip_hdr = output_pkt + 14;
    memset(ip_hdr, 0, 20);
    ip_hdr[0] = 0x45; /* IPv4, IHL=5 */
    *(uint16_t *)(ip_hdr + 2) = htons(20 + 8 + 8 + input_len); /* Total length */
    ip_hdr[8] = 64; /* TTL */
    ip_hdr[9] = 17; /* UDP protocol */
    *(uint32_t *)(ip_hdr + 12) = tunnel->local_ip;
    *(uint32_t *)(ip_hdr + 16) = tunnel->remote_ip;

    /* Outer UDP header */
    uint8_t *udp_hdr = output_pkt + 14 + 20;
    *(uint16_t *)(udp_hdr + 0) = htons(4789); /* VXLAN port */
    *(uint16_t *)(udp_hdr + 2) = htons(4789);
    *(uint16_t *)(udp_hdr + 4) = htons(8 + 8 + input_len); /* UDP length */
    *(uint16_t *)(udp_hdr + 6) = 0; /* Checksum optional for UDP */

    /* Outer Ethernet header */
    memset(output_pkt, 0xFF, 6); /* Dest MAC (broadcast for now) */
    memset(output_pkt + 6, 0x00, 6); /* Src MAC */
    *(uint16_t *)(output_pkt + 12) = htons(0x0800); /* IPv4 ethertype */

    *output_len = header_size + input_len;
    state->packets_encapsulated++;
    return true;
}

bool vxlan_decapsulate(vxlan_state_t *state, const uint8_t *vxlan_pkt, uint32_t pkt_len,
                       uint8_t *output_pkt, uint32_t *output_len) {
    if (pkt_len < 50) {
        return false; /* Too small */
    }

    /* Skip outer headers: Eth(14) + IP(20) + UDP(8) + VXLAN(8) = 50 */
    uint32_t inner_start = 50;
    *output_len = pkt_len - inner_start;

    memcpy(output_pkt, vxlan_pkt + inner_start, *output_len);
    state->packets_decapsulated++;
    return true;
}

void vxlan_get_stats(vxlan_state_t *state, uint64_t *encap, uint64_t *decap) {
    *encap = state->packets_encapsulated;
    *decap = state->packets_decapsulated;
}
