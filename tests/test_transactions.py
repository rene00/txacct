import json

def test_transactions(client):
    with client:
        path: str ="/transactions/"
        response=client.post(path,
         json={"memo":"foo"},
        )
        data = json.loads(response.text)
        assert response.status_code == 200
        assert data.get("memo") == "foo"
