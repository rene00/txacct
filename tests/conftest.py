import pytest
from flask import Flask
from flask_migrate import upgrade
from txacct import create_app
from txacct.importer import Importer
from txacct.model import db as _db

postcodes = [
    {
        "postcode": "3000",
        "locality": "MELBOURNE",
        "state": "VIC",
        "sa3": "20604",
        "sa3name": "Melbourne City",
        "sa4": "206",
        "sa4name": "Melbourne - Inner",
    },
    {
        "postcode": "3006",
        "locality": "SOUTH WHARF",
        "state": "VIC",
        "sa3": "20605",
        "sa3name": "Port Phillip",
        "sa4": "206",
        "sa4name": "Melbourne - Inner",
    },
    {
        "postcode": "2000",
        "locality": "BARANGAROO",
        "state": "NSW",
        "sa3": "11703",
        "sa3name": "Sydney Inner City",
        "sa4": "117",
        "sa4name": "Sydney - City and Inner South",
    },
]

from pytest_postgresql.janitor import DatabaseJanitor

TEST_DB = "txacct"


@pytest.fixture
def connection(postgresql_proc):
    USER = postgresql_proc.user
    HOST = postgresql_proc.host
    PORT = postgresql_proc.port
    with DatabaseJanitor(USER, HOST, PORT, TEST_DB, 12.2):
        yield f"postgresql://{USER}:@{HOST}:{PORT}/{TEST_DB}"


@pytest.fixture
def app(connection):
    app: Flask = create_app(dict(SQLALCHEMY_DATABASE_URI=connection))

    with app.app_context():
        upgrade()
        Importer(_db.session, postcode_data=postcodes).do()

    yield app


@pytest.fixture
def client(app):
    return app.test_client()
