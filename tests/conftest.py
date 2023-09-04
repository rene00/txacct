import pytest
import tempfile
from flask import Flask
from flask_migrate import upgrade
from txacct import create_app
import os


@pytest.fixture
def app():
    db_fd, db_path = tempfile.mkstemp()
    app: Flask = create_app(
        dict(
            SQLALCHEMY_DATABASE_URI="sqlite:///{0}".format(db_path),
        )
    )

    with app.app_context():
        upgrade()

    yield app
    os.close(db_fd)
    os.unlink(db_path)


@pytest.fixture
def client(app):
    return app.test_client()
