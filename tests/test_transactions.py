import json
from collections import namedtuple

# Example response:
# {
#  "id": 1,
#  "locality": {
#    "name": "MELBOURNE",
#    "postcode": "3000",
#    "state": "VIC"
#    }
#  },
#  "memo": "TEST MELBOURNE VIC"
# }


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
            assert data.get("locality").get("state") == i.state


def test_transactions_locality(client):
    testcase = namedtuple("testcase", ["memo", "name", "postcode", "state"])
    testcases = (
        testcase(
            memo="TEST MELBOURNE VIC",
            name="MELBOURNE",
            postcode="3000",
            state="VIC",
        ),
        testcase(
            memo="TEST SOUTH WHARF VIC",
            name="SOUTH WHARF",
            postcode="3006",
            state="VIC",
        ),
        testcase(
            memo="TEST SOUTH WHARFVIC",
            name="SOUTH WHARF",
            postcode="3006",
            state="VIC",
        ),
        testcase(
            memo="TEST SOUTH WHARF VI",
            name="SOUTH WHARF",
            postcode="3006",
            state="VIC",
        ),
    )

    with client:
        for i in testcases:
            response = client.post(
                "/transactions/",
                json={"memo": i.memo},
            )
            data = json.loads(response.text)
            print(data)
            assert response.status_code == 200
            assert data["locality"]["name"] == i.name
            assert data["locality"]["postcode"] == i.postcode
            assert data["locality"]["state"] == i.state
