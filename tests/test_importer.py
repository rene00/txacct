from txacct.model import State, StatisticalArea3, StatisticalArea4, Postcode
from txacct.model import db
from sqlalchemy import select
from collections import namedtuple


def test_importer_state(app):
    with app.app_context():
        for i in ("VIC", "NSW"):
            assert db.session.execute(select(State).where(State.name == i)).one()


def test_importer_sa3(app):
    with app.app_context():
        for i in ("Melbourne City", "Port Phillip", "Sydney Inner City"):
            assert db.session.execute(
                select(StatisticalArea3).where(StatisticalArea3.name == i)
            ).one()


def test_importer_sa4(app):
    with app.app_context():
        for i in ("Melbourne - Inner", "Sydney - City and Inner South"):
            assert db.session.execute(
                select(StatisticalArea4).where(StatisticalArea4.name == i)
            ).one()


def test_importer_postcode(app):
    postcode = namedtuple("Postcode", ["postcode", "locality"])
    with app.app_context():
        for i in (
            postcode("3006", "SOUTH WHARF"),
            postcode("3000", "MELBOURNE"),
            postcode("2000", "BARANGAROO"),
        ):
            assert db.session.execute(
                select(Postcode).where(
                    Postcode.postcode == i.postcode, Postcode.locality == i.locality
                )
            ).one()
