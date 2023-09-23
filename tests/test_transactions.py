import json
from collections import namedtuple


def test_transactions(client):
    with client:
        path: str = "/transactions/"
        response = client.post(
            path,
            json={"memo": "foo"},
        )
        data = json.loads(response.text)
        assert response.status_code == 200
        assert data.get("memo") == "foo"


def test_transactions_state(client):
    testcase = namedtuple("testcase", ["memo", "state"])
    testcases = (
        testcase("TEST MELBOURNE VIC", "VIC"),
        testcase("TEST SOUTH WHARF VIC", "VIC"),
        testcase("TEST SOUTH WHARFVIC", "VIC"),
    )

    with client:
        for i in testcases:
            response = client.post(
                "/transactions/",
                json={"memo": i.memo},
            )
            data = json.loads(response.text)
            assert response.status_code == 200
            assert data.get("locality").get("state").get("name") == i.state
