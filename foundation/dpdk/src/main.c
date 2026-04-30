#include <stdio.h>
#include <string.h>
#include "forwarding.h"
#include "vxlan.h"

int main(int argc, char *argv[]) {
    (void)argc;
    (void)argv;

    printf("DPDK Handler Starting\n");

    /* Initialize forwarding */
    forwarding_state_t fwd_state;
    forwarding_init(&fwd_state);
    forwarding_add_route(&fwd_state, 0x0A000100, 0xFFFFFF00, 0x0A000001, 1);

    /* Initialize VXLAN */
    vxlan_state_t vxlan_state;
    vxlan_init(&vxlan_state);
    vxlan_add_tunnel(&vxlan_state, 1, 0xC0A80101, 0xC0A80102, 100); /* 192.168.1.1 -> .2, VNI 100 */

    /* Test packet - 60 bytes total (IPv4 header + 40-byte payload) */
    uint8_t pkt_data[60];
    ipv4_hdr_t *pkt = (ipv4_hdr_t *)pkt_data;
    memset(pkt_data, 0, 60);
    pkt->version_ihl = 0x45;
    pkt->dscp_ecn = 0x00;
    pkt->total_length = 60;
    pkt->identification = 1;
    pkt->flags_frag_offset = 0x4000;
    pkt->ttl = 64;
    pkt->protocol = 6;
    pkt->checksum = 0;
    pkt->src_ip = 0x0A000102;  /* 10.0.1.2 */
    pkt->dest_ip = 0x0A000103; /* 10.0.1.3 */

    /* Forward decision */
    int egress = forwarding_decide(&fwd_state, pkt);
    printf("Forwarding decision: egress_port=%d\n", egress);

    /* Encapsulate in VXLAN */
    uint8_t encap_pkt[256];
    uint32_t encap_len;
    bool encap_ok = vxlan_encapsulate(&vxlan_state, 1, pkt_data, 60,
                                      encap_pkt, &encap_len);
    printf("VXLAN encapsulation: %s, len=%u\n", encap_ok ? "OK" : "FAIL", encap_len);

    /* Decapsulate */
    uint8_t decap_pkt[256];
    uint32_t decap_len;
    bool decap_ok = vxlan_decapsulate(&vxlan_state, encap_pkt, encap_len,
                                      decap_pkt, &decap_len);
    printf("VXLAN decapsulation: %s, len=%u\n", decap_ok ? "OK" : "FAIL", decap_len);

    /* Stats */
    uint64_t fwd_count, drop_count;
    forwarding_get_stats(&fwd_state, &fwd_count, &drop_count);
    printf("Forwarding stats: fwd=%lu, drop=%lu\n", fwd_count, drop_count);

    uint64_t encap_count, decap_count;
    vxlan_get_stats(&vxlan_state, &encap_count, &decap_count);
    printf("VXLAN stats: encap=%lu, decap=%lu\n", encap_count, decap_count);

    return 0;
}
