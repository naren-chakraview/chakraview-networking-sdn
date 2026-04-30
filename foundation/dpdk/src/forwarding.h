#ifndef FORWARDING_H
#define FORWARDING_H

#include <stdint.h>
#include <stdbool.h>

/* Maximum number of routing table entries */
#define MAX_ROUTES 1000
#define MAX_FILTERS 500

/* Packet header structures */
typedef struct {
    uint8_t dest_mac[6];
    uint8_t src_mac[6];
    uint16_t ethertype;
} eth_hdr_t;

typedef struct {
    uint8_t version_ihl;
    uint8_t dscp_ecn;
    uint16_t total_length;
    uint16_t identification;
    uint16_t flags_frag_offset;
    uint8_t ttl;
    uint8_t protocol;
    uint16_t checksum;
    uint32_t src_ip;
    uint32_t dest_ip;
} ipv4_hdr_t;

/* Routing table entry */
typedef struct {
    uint32_t dest_ip;
    uint32_t mask;
    uint32_t next_hop;
    uint16_t egress_port;
} route_entry_t;

/* Filter rule (packet filter) */
typedef struct {
    uint32_t src_ip;
    uint32_t dest_ip;
    uint8_t action; /* 0 = drop, 1 = forward */
} filter_rule_t;

/* Forwarding engine state */
typedef struct {
    route_entry_t routes[MAX_ROUTES];
    int route_count;
    filter_rule_t filters[MAX_FILTERS];
    int filter_count;
    uint64_t packets_forwarded;
    uint64_t packets_dropped;
} forwarding_state_t;

/* Initialize forwarding engine */
void forwarding_init(forwarding_state_t *state);

/* Add route to routing table */
bool forwarding_add_route(forwarding_state_t *state, uint32_t dest_ip,
                          uint32_t mask, uint32_t next_hop, uint16_t egress_port);

/* Add filter rule */
bool forwarding_add_filter(forwarding_state_t *state, uint32_t src_ip,
                           uint32_t dest_ip, uint8_t action);

/* Forward a packet: returns next_hop port, -1 if drop */
int forwarding_decide(forwarding_state_t *state, const ipv4_hdr_t *pkt);

/* Get statistics */
void forwarding_get_stats(forwarding_state_t *state, uint64_t *pkt_fwd, uint64_t *pkt_drop);

#endif
