import json
from collections import namedtuple

# Example response:
# {
#  "id": 1,
#  "locality": {
#    "names": {
#      "100": {
#        "name": "MELBOURNE",
#        "postcode": "3000",
#        "sa3": {
#          "name": "Melbourne City"
#        },
#        "sa4": {
#          "name": "Melbourne - Inner"
#        }
#      }
#    },
#    "state": {
#      "name": "VIC"
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
            assert data.get("locality").get("state").get("name") == i.state


def test_transactions_locality(client):
    testcase = namedtuple("testcase", ["memo", "name", "postcode", "sa3", "sa4"])
    testcases = (
        testcase(
            memo="TEST MELBOURNE VIC",
            name="MELBOURNE",
            postcode="3000",
            sa3="Melbourne City",
            sa4="Melbourne - Inner",
        ),
        testcase(
            memo="TEST MELBOURNEVIC",
            name="MELBOURNE",
            postcode="3000",
            sa3="Melbourne City",
            sa4="Melbourne - Inner",
        ),
        testcase(
            memo="TEST SOUTH WHARF VIC",
            name="SOUTH WHARF",
            postcode="3006",
            sa3="Port Phillip",
            sa4="Melbourne - Inner",
        ),
        testcase(
            memo="TEST SOUTH WHARFVIC",
            name="SOUTH WHARF",
            postcode="3006",
            sa3="Port Phillip",
            sa4="Melbourne - Inner",
        ),
    )

    with client:
        for i in testcases:
            response = client.post(
                "/transactions/",
                json={"memo": i.memo},
            )
            data = json.loads(response.text)
            assert response.status_code == 200
            assert data["locality"]["names"]["100"]["name"] == i.name
            assert data["locality"]["names"]["100"]["postcode"] == i.postcode
            assert data["locality"]["names"]["100"]["sa3"]["name"] == i.sa3
            assert data["locality"]["names"]["100"]["sa4"]["name"] == i.sa4
