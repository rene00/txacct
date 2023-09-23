from flask.cli import AppGroup, with_appcontext
from flask import current_app
import nltk
import requests
import click
import json
from txacct.model import db
from .model import StatisticalArea3, StatisticalArea4, State, Postcode
from sqlalchemy.exc import IntegrityError

postcode_cli = AppGroup("postcode")
nltk_cli = AppGroup("nltk")


@nltk_cli.command("download")
def nltk_download():
    nltk.download("punkt")


@postcode_cli.command("import")
@with_appcontext
def postcode_import():
    config = current_app.config
    postcode_url: str = config.get("POSTCODE_URL", "")
    resp = requests.get(postcode_url)
    if resp.status_code != 200:
        raise click.ClickException("failed to download")

    for i in json.loads(resp.text):
        sa3 = None
        if i.get("sa3") != "" and i.get("sa3name") != "":
            sa3 = StatisticalArea3(
                code=i.get("sa3"),
                name=i.get("sa3name"),
            )
            db.session.add(sa3)
            try:
                db.session.commit()
            except IntegrityError:
                db.session.rollback()
                sa3 = (
                    db.session.query(StatisticalArea3)
                    .filter_by(code=int(sa3.code), name=i.get("sa3name"))
                    .one()
                )

        sa4 = None
        if i.get("sa4") != "" and i.get("sa4name") != "":
            sa4 = StatisticalArea4(
                code=i.get("sa4"),
                name=i.get("sa4name"),
            )
            db.session.add(sa4)
            try:
                db.session.commit()
            except IntegrityError:
                db.session.rollback()
                sa4 = (
                    db.session.query(StatisticalArea4)
                    .filter_by(code=int(sa4.code), name=sa4.name)
                    .one()
                )

        state = None
        if i.get("state") != "":
            state = State(name=i.get("state"))
            db.session.add(state)
            try:
                db.session.commit()
            except IntegrityError:
                db.session.rollback()
                state = db.session.query(State).filter_by(name=state.name).one()

        postcode: Postcode = Postcode(
            postcode=i.get("postcode"),
            locality=i.get("locality"),
            state=state,
            sa3=sa3,
            sa4=sa4,
        )

        db.session.add(postcode)
        try:
            db.session.commit()
        except IntegrityError:
            db.session.rollback()

    return None
