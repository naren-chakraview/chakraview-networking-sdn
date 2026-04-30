"""End-to-end integration tests"""

import requests
import time
import pytest

CONTROLLER_URL = "http://localhost:8080"


class TestTopology:
    def test_health_check(self):
        """Verify controller health"""
        response = requests.get(f"{CONTROLLER_URL}/api/v1/health")
        assert response.status_code == 200
        data = response.json()
        assert data['status'] == 'healthy'

    def test_topology_summary(self):
        """Get topology summary"""
        response = requests.get(f"{CONTROLLER_URL}/api/v1/topology")
        assert response.status_code == 200
        data = response.json()
        assert 'summary' in data
        assert data['status'] == 'ok'

    def test_device_listing(self):
        """List registered devices"""
        response = requests.get(f"{CONTROLLER_URL}/api/v1/topology/devices")
        assert response.status_code == 200
        data = response.json()
        assert 'devices' in data


if __name__ == '__main__':
    pytest.main([__file__, '-v'])
