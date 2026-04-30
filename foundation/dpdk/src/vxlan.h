#ifndef VXLAN_H
#define VXLAN_H

#include <stdint.h>
#include <stdbool.h>

#define MAX_VXLAN_TUNNELS 100

/* VXLAN header (8 bytes after UDP) */
typedef struct {
    uint8_t flags;
    uint8_t reserved1[3];
    uint32_t vni; /* 24-bit VNI + 8 reserved */
} vxlan_hdr_t;

/* VXLAN tunnel definition */
typedef struct {
    uint32_t tunnel_id;
    uint32_t local_ip;
    uint32_t remote_ip;
    uint32_t vni;
    bool active;
} vxlan_tunnel_t;

/* VXLAN state */
typedef struct {
    vxlan_tunnel_t tunnels[MAX_VXLAN_TUNNELS];
    int tunnel_count;
    uint64_t packets_encapsulated;
    uint64_t packets_decapsulated;
} vxlan_state_t;

/* Initialize VXLAN state */
void vxlan_init(vxlan_state_t *state);

/* Add a VXLAN tunnel */
bool vxlan_add_tunnel(vxlan_state_t *state, uint32_t tunnel_id, uint32_t local_ip,
                      uint32_t remote_ip, uint32_t vni);

/* Encapsulate packet in VXLAN (output buffer must be at least input_len + 50 bytes) */
bool vxlan_encapsulate(vxlan_state_t *state, uint32_t tunnel_id,
                       const uint8_t *input_pkt, uint32_t input_len,
                       uint8_t *output_pkt, uint32_t *output_len);

/* Decapsulate VXLAN packet */
bool vxlan_decapsulate(vxlan_state_t *state, const uint8_t *vxlan_pkt, uint32_t pkt_len,
                       uint8_t *output_pkt, uint32_t output_buf_size, uint32_t *output_len);

/* Get statistics */
void vxlan_get_stats(vxlan_state_t *state, uint64_t *encap, uint64_t *decap);

#endif
