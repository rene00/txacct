from flask import Flask, Blueprint, jsonify, Response, request, abort, current_app
from flask.cli import AppGroup, with_appcontext
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from typing import Any, List
from sqlalchemy import func, String, ForeignKey, Integer
from sqlalchemy.orm import mapped_column, Mapped, relationship
from sqlalchemy.exc import IntegrityError, NoResultFound
import sys
import requests
import json
import nltk
from nltk.tokenize import word_tokenize

db: SQLAlchemy = SQLAlchemy()

transactions_bp: Blueprint = Blueprint("transactions", __name__, url_prefix="/transactions") 


class Transaction(db.Model):
    id: Mapped[int] = mapped_column(primary_key=True)
    memo: Mapped[str] = mapped_column(String(), nullable=False)

class Postcode(db.Model):
    __table_args__: tuple[Any] = (db.UniqueConstraint("postcode", "locality", name="postcode_locality"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    postcode: Mapped[str] = mapped_column(String(), nullable=False)
    locality: Mapped[str] = mapped_column(String(), nullable=False)
    state_id: Mapped[int] = mapped_column(ForeignKey("state.id"))
    state: Mapped["State"] = relationship(back_populates="postcodes")
    sa3_id: Mapped[int] = mapped_column(ForeignKey("sa3.id"), nullable=True)
    sa3: Mapped["StatisticalArea3"] = relationship(back_populates="postcodes")
    sa4_id: Mapped[int] = mapped_column(ForeignKey("sa4.id"), nullable=True)
    sa4: Mapped["StatisticalArea4"] = relationship(back_populates="postcodes")

class StatisticalArea3(db.Model):
    __tablename__: str = "sa3"
    __table_args__: tuple[Any] = (db.UniqueConstraint("code", "name"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    code: Mapped[int] = mapped_column(Integer(), nullable=False)
    name: Mapped[str] = mapped_column(String(), nullable=False)
    postcodes: Mapped[List["Postcode"]] = relationship(back_populates="sa3")

class StatisticalArea4(db.Model):
    __tablename__: str = "sa4"
    __table_args__: tuple[Any] = (db.UniqueConstraint("code", "name"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    code: Mapped[int] = mapped_column(Integer(), nullable=False)
    name: Mapped[str] = mapped_column(String(), nullable=False)
    postcodes: Mapped[List["Postcode"]] = relationship(back_populates="sa4")

class State(db.Model):
    id: Any = db.Column(db.Integer, primary_key=True)
    name: Any = db.Column(db.String, nullable=False, unique=True)
    postcodes: Mapped[List["Postcode"]] = relationship(back_populates="state")

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

    tm = TransactionMeta(memo=transaction.memo, db=db)

    d = { 
        "id": transaction.id,
        "memo": transaction.memo,
    }

    locality = dict()
    state = tm.state()
    if state is not None:
        locality["state"] = dict({
            "name": state.name,
        })

        postcode = tm.postcode(locality=tm.tokenized[-2], state=state)
        print(postcode)
        if len(postcode) >= 1:
            names = []
            for i in postcode:
                d2 = dict({"name": i.locality, "postcode": i.postcode})
                if i.sa3:
                    d2["sa3"] = dict({
                        "name": i.sa3.name,
                    })
                if i.sa4:
                    d2["sa4"] = dict({
                        "name": i.sa4.name,
                    })
                names.append(d2)

            locality["names"] = names
        

    d["locality"] = locality

    return jsonify(d)

class TransactionMeta:
    def __init__(self, memo: str, db) -> None:
        self.memo: str = memo
        self.db = db
        self.tokenized: List[str] = word_tokenize(self.memo)

    """ state returns the Australian state of the transaction. """
    def state(self) -> (State | None):
        # Take last item from tokenized and see if its a match on state.
        l: str = self.tokenized[-1]
        try:
            s = self.db.session.query(State).filter(func.lower(State.name) == l.lower()).one()
        except NoResultFound:
            s = None
        return s

    def postcode(self, locality: str, state: State) -> List[Postcode]:
        return self.db.session.query(Postcode).filter(Postcode.locality == locality, Postcode.state_id == state.id).all()

postcode_cli = AppGroup('postcode')
nltk_cli = AppGroup('nltk')

@nltk_cli.command('download')
def nltk_download():
    nltk.download('punkt')

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

def create_app(test_config: dict | None=None) -> Flask:
    app: Flask = Flask(__name__)
    if test_config:
        app.config.from_mapping(test_config)
    else:
        app.config.from_prefixed_env(prefix="TXACCT")
    db.init_app(app)
    Migrate(app, db)
    app.register_blueprint(transactions_bp)
    app.cli.add_command(postcode_cli)
    app.cli.add_command(nltk_cli)
    return app
