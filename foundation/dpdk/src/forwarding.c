#include "forwarding.h"
#include <string.h>
#include <stdio.h>

void forwarding_init(forwarding_state_t *state) {
    memset(state, 0, sizeof(forwarding_state_t));
    state->packets_forwarded = 0;
    state->packets_dropped = 0;
}

bool forwarding_add_route(forwarding_state_t *state, uint32_t dest_ip,
                          uint32_t mask, uint32_t next_hop, uint16_t egress_port) {
    if (state->route_count >= MAX_ROUTES) {
        return false;
    }

    route_entry_t *route = &state->routes[state->route_count];
    route->dest_ip = dest_ip;
    route->mask = mask;
    route->next_hop = next_hop;
    route->egress_port = egress_port;

    state->route_count++;
    return true;
}

bool forwarding_add_filter(forwarding_state_t *state, uint32_t src_ip,
                           uint32_t dest_ip, uint8_t action) {
    if (state->filter_count >= MAX_FILTERS) {
        return false;
    }

    filter_rule_t *filter = &state->filters[state->filter_count];
    filter->src_ip = src_ip;
    filter->dest_ip = dest_ip;
    filter->action = action;

    state->filter_count++;
    return true;
}

/* Check if packet passes filters (1 = pass, 0 = drop) */
static int filter_check(forwarding_state_t *state, const ipv4_hdr_t *pkt) {
    for (int i = 0; i < state->filter_count; i++) {
        filter_rule_t *f = &state->filters[i];
        /* Exact match for now; TODO: EXTEND: add CIDR matching */
        if (f->src_ip == pkt->src_ip && f->dest_ip == pkt->dest_ip) {
            return f->action;
        }
    }
    return 1; /* Default allow */
}

/* Find longest matching prefix route */
static int route_lookup(forwarding_state_t *state, uint32_t dest_ip) {
    int best_match = -1;
    int best_prefix_len = -1;

    for (int i = 0; i < state->route_count; i++) {
        route_entry_t *r = &state->routes[i];
        if ((dest_ip & r->mask) == (r->dest_ip & r->mask)) {
            /* Count prefix length (number of 1s in mask) */
            int prefix_len = __builtin_popcount(r->mask);
            if (prefix_len > best_prefix_len) {
                best_prefix_len = prefix_len;
                best_match = i;
            }
        }
    }

    return best_match;
}

int forwarding_decide(forwarding_state_t *state, const ipv4_hdr_t *pkt) {
    /* Check filters first */
    if (!filter_check(state, pkt)) {
        state->packets_dropped++;
        return -1;
    }

    /* Lookup route */
    int route_idx = route_lookup(state, pkt->dest_ip);
    if (route_idx == -1) {
        state->packets_dropped++;
        return -1; /* No route */
    }

    state->packets_forwarded++;
    return state->routes[route_idx].egress_port;
}

void forwarding_get_stats(forwarding_state_t *state, uint64_t *pkt_fwd, uint64_t *pkt_drop) {
    *pkt_fwd = state->packets_forwarded;
    *pkt_drop = state->packets_dropped;
}
