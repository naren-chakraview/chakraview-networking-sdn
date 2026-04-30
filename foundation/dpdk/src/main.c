#include <stdio.h>
#include "forwarding.h"

int main(int argc, char *argv[]) {
    (void)argc;
    (void)argv;

    printf("DPDK Handler Starting\n");

    forwarding_state_t state;
    forwarding_init(&state);

    /* Example: add a route */
    forwarding_add_route(&state, 0x0A000100, 0xFFFFFF00, 0x0A000001, 1);

    /* Example: create a test packet */
    ipv4_hdr_t pkt = {
        .version_ihl = 0x45,
        .dscp_ecn = 0x00,
        .total_length = 60,
        .identification = 1,
        .flags_frag_offset = 0x4000,
        .ttl = 64,
        .protocol = 6, /* TCP */
        .checksum = 0,
        .src_ip = 0x0A000102,  /* 10.0.1.2 */
        .dest_ip = 0x0A000103  /* 10.0.1.3 */
    };

    int egress = forwarding_decide(&state, &pkt);
    printf("Packet forwarding decision: egress_port=%d\n", egress);

    uint64_t fwd, drop;
    forwarding_get_stats(&state, &fwd, &drop);
    printf("Stats: forwarded=%lu, dropped=%lu\n", fwd, drop);

    return 0;
}
