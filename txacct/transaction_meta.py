from .model import ANZSIC, Organisation, State, Postcode, BusinessCode
from typing import List
from sqlalchemy import func, select, or_, and_
from txacct.model import db

import spacy
from spacy.language import Language
from spacy.matcher import Matcher
from spacy.tokens import Doc, Span, Token
import re


@Language.factory("postcode")
class PostcodeComponent:
    def __init__(self, nlp, name, label):
        self.nlp = nlp
        self.label = label
        self.name = name

        Token.set_extension("is_locality", default=None, force=True)
        Doc.set_extension("postcode", default=None, force=True)

        self.db = db

    def __is_locality(self, token: Token, doc: Doc) -> bool:
        """
        __is_locality() checks the doc to see if a token within the doc
        which shares the same text as the token arg already has
        token._.is_locality. If a match is found and is_locality is set, return
        True else return False.

        Only one token within a doc which shares the same text as other tokens
        should have is_locality set and that token should be the higher indexed
        token.

        example memo:

            the haus hanhndorf hanhndorf sa

        The last "hanhndorf" token should have is_locality set and not the first.
        """
        for t in doc:
            if t.text.lower() == token.text.lower() and t._.is_locality:
                return True
        return False

    def __call__(self, doc):
        for idx in range(len(doc)):
            if doc._.get("postcode") is not None:
                break

            token = doc[len(doc) - (idx + 1)]

            rows = self.db.session.execute(
                select(Postcode).where(
                    func.lower(Postcode.locality) == token.text.lower()
                )
            ).all()

            state = doc._.state

            for row in rows:
                postcode = row[0]
                if state is not None and postcode.state != state:
                    continue
                doc._.set("postcode", postcode)
                if not self.__is_locality(token, doc):
                    token._.set("is_locality", True)
                    break

            # If postcode still not set, attempt to find a match combining
            # the previous and current token in the loop.
            if doc._.postcode is None and idx >= 1:
                token_prev = doc[len(doc) - (idx + 2)]
                rows = self.db.session.execute(
                    select(Postcode).where(
                        func.lower(Postcode.locality)
                        == f"{token_prev.text} {token.text}".lower()
                    )
                ).all()

                for row in rows:
                    postcode = row[0]
                    token._.set("is_locality", True)
                    token_prev._.set("is_locality", True)
                    doc._.set("postcode", postcode)
                    break

        # If postcode_name still not set, take last token and strip state name
        # from it. If the strip was successful, take this newly stripped text
        # and previous token text and then regexp on all postnames.
        if doc._.postcode is None:
            last_token = doc[-1]
            rows = self.db.session.execute(select(State)).all()
            for row in rows:
                state = row[0]
                if doc._.postcode is not None:
                    break
                if state.name.lower() in last_token.text.lower() and len(
                    state.name
                ) < len(last_token.text):
                    s = last_token.text.lower().strip(state.name.lower())
                    second_last_token = doc[len(doc) - 2]
                    search = f"{second_last_token.text} {s}%".lower()
                    rows = self.db.session.execute(
                        select(Postcode).where(Postcode.locality.ilike(search))
                    ).all()
                    postcode = None
                    for r in rows:
                        postcode = r[0]
                    if postcode is None:
                        continue
                    doc._.set("postcode", postcode)
                    last_token._.set("is_locality", True)
                    second_last_token._.set("is_locality", True)
                    break

        return doc

    def has_postcode(self, doc) -> bool:
        if doc._.postcode_name is None:
            return False
        return True


@Language.factory("state")
class StateComponent:
    def __init__(self, nlp, name, label):
        self.nlp = nlp
        self.name = name
        self.label = label

        Token.set_extension("is_locality", default=None, force=True)
        Doc.set_extension("state", default=None, force=True)

        self.db = db

    def __call__(self, doc):
        states = self.db.session.execute(select(State)).all()
        for token in reversed(doc):
            if doc._.state is not None:
                break
            for row in states:
                state = row[0]
                regexp = rf"^(?P<prefix>\w+)?(?P<state>{state.name}|{state.name[:2]})$"
                if re.match(regexp, token.text):
                    doc._.set("state", state)
                    token._.set("is_locality", True)
                    break

        if not doc._.state:
            # Lookup state using postcode. Get last 2 tokens from doc and
            # lookup in Postcode. If Postcode found, use state from that.
            tokens = (doc[-1], doc[len(doc) - 2])

            for token in tokens:
                rows = self.db.session.execute(
                    select(Postcode).where(Postcode.locality == token.text)
                ).all()

                if len(rows) == 0:
                    continue

                # Multiple postcodes can be found with the above query. Select
                # the first postcode. This can be improved though for now it's
                # good enough.
                state = rows[0]
                doc._.set("state", state.name)
                token._.set("is_locality", True)
                break

        return doc


class TransactionMeta:
    def __init__(self, memo: str, db) -> None:
        self.memo: str = memo
        self.db = db

        self.nlp = spacy.blank("en")
        self.nlp.add_pipe("state", config={"label": "STATE"})
        self.nlp.add_pipe("postcode", config={"label": "POSTCODE"})
        self.doc = self.nlp(memo)

        self.__cache = {}

    def state(self) -> State | None:
        """state returns the Australian state of the transaction."""
        return self.doc._.state

    def postcode(self) -> Postcode | None:
        return self.doc._.postcode

    def __organisation_query(self) -> list:
        """
        Return a list of WHERE ILIKE clauses for Organisation name.

        Potential Improvements:
          - Split token.text if contains CamelCase (i.e. MelbournePool)
        """
        queries = []
        for token in self.doc:
            if token._.is_locality:
                continue
            if len(queries) == 0:
                query = token.text
            else:
                query = queries[-1] + f" {token}"
            queries.append(query)
        return queries

    def organisation(self) -> Organisation | None:
        cache = self.__cache.get("organisation")
        if cache is not None:
            return cache

        postcode = self.postcode()

        ret = None
        queries = self.__organisation_query()
        for query in reversed(queries):
            search = f"{query}%".lower()
            if postcode is None:
                stmt = (
                    select(Organisation, BusinessCode)
                    .where(Organisation.name.ilike(search))
                    .join_from(Organisation, BusinessCode)
                )
            else:
                stmt = (
                    select(Organisation, BusinessCode, Postcode)
                    .join_from(Organisation, BusinessCode)
                    .join_from(Organisation, Postcode)
                    .where(Organisation.name.ilike(search))
                    .where(Postcode.id == postcode.id)
                )

            # IMPROVEMENT: instead of returning first, return all() and provide
            # suggested results to client.
            ret = self.db.session.scalars(stmt).first()
            if ret is not None:
                break

        self.__cache["organisation"] = ret

        return ret

    def address(self) -> str | None:
        organisation = self.organisation()
        if organisation is None:
            return None
        address = organisation.address

        if organisation.postcode:
            address = f"{address}, {organisation.postcode.locality}, {organisation.postcode.state.name}"

        return address

    def business_code(self) -> BusinessCode | None:
        o = self.organisation()
        if o is None:
            return None
        return o.business_code

    def anzsic(self):
        o = self.organisation()
        if o is None:
            return None
        return o.anzsic
