from flask import Flask, Blueprint, jsonify, Response, request, abort, current_app
from flask.cli import AppGroup, with_appcontext
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from typing import Any
import click
from sqlalchemy.orm import RelationshipProperty
from sqlalchemy.exc import IntegrityError
import sys
import requests
import json
import pprint
from marshmallow_sqlalchemy import SQLAlchemySchema, auto_field

db: SQLAlchemy = SQLAlchemy()

transactions_bp: Blueprint = Blueprint("transactions", __name__, url_prefix="/transactions") 

class Transaction(db.Model):
    id: Any = db.Column(db.Integer, primary_key=True)
    memo: Any = db.Column(db.String, nullable=False)

class TransactionSchema(SQLAlchemySchema):
    class Meta:
        model = Transaction
        include_relationships: bool = True
        load_instance: bool = True

    id = auto_field()
    memo = auto_field()

class Postcode(db.Model):
    __table_args__: tuple[Any] = (db.UniqueConstraint("postcode", "locality"),)

    id: Any = db.Column(db.Integer, primary_key=True)
    postcode: Any = db.Column(db.String, nullable=False)
    locality: Any = db.Column(db.String, nullable=False)
    state: Any = db.Column(db.String, nullable=False)
    sa3_id: Any = db.Column(db.Integer, db.ForeignKey("sa3.id"))
    sa4_id: Any = db.Column(db.Integer, db.ForeignKey("sa4.id"))

class StatisticalArea3(db.Model):
    __tablename__: str = "sa3"
    __table_args__: tuple[Any] = (db.UniqueConstraint("code", "name"),)

    id: Any = db.Column(db.Integer, primary_key=True)
    code: Any = db.Column(db.Integer, nullable=False)
    name: Any = db.Column(db.String, nullable=False)
    postcodes = db.relationship("Postcode", backref="sa3")

class StatisticalArea4(db.Model):
    __tablename__: str = "sa4"
    __table_args__: tuple[Any] = (db.UniqueConstraint("code", "name"),)

    id: Any = db.Column(db.Integer, primary_key=True)
    code: Any = db.Column(db.Integer, nullable=False)
    name: Any = db.Column(db.String, nullable=False)
    postcodes = db.relationship("Postcode", backref="sa4")

@transactions_bp.route("/", methods=["POST"])
def transactions() -> Response:
    if not request.is_json:
        abort(403)
    data: Any = request.json
    transaction: Transaction = Transaction(
            memo=data.get('memo')
    )
    db.session.add(transaction)
    db.session.commit()

    return jsonify(TransactionSchema().dump(transaction))

postcode_cli = AppGroup('postcode')
@postcode_cli.command('import')
@with_appcontext
def postcode_import():
    config = current_app.config
    postcode_url = config.get("POSTCODE_URL")
    debug = config.get("DEBUG", False)
    resp = requests.get(postcode_url)
    if resp.status_code != 200:
        sys.exit("failed to download")

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
                sa3 = db.session.query(StatisticalArea3).filter_by(code=int(sa3.code), name=i.get("sa3name")).one()

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
                sa4 = db.session.query(StatisticalArea4).filter_by(code=int(sa4.code), name=sa4.name).one()

        postcode: Postcode = Postcode(
            postcode=i.get("postcode"),
            locality=i.get("locality"),
            state=i.get("state"),
            sa3=sa3,
            sa4=sa4,
        )

        db.session.add(postcode)
        try:
            db.session.commit()
        except IntegrityError:
            db.session.rollback()

    return None

def create_app() -> Flask:
    app: Flask = Flask(__name__)
    app.config.from_prefixed_env(prefix="TXACCT")
    db.init_app(app)
    Migrate(app, db)
    app.register_blueprint(transactions_bp)
    app.cli.add_command(postcode_cli)
    return app
