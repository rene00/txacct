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

        self.__postcode: Postcode | None = None
        self.__organisation: Organisation | None = None

    def state(self) -> State | None:
        """state returns the Australian state of the transaction."""
        return self.doc._.state

    def postcode(self) -> Postcode | None:
        if self.__postcode is not None:
            return Postcode
        return self.doc._.postcode

    def __query_skip_token(self, token: Token) -> bool:
        if token._.is_locality:
            return True

        # skip key is token index. The value is the text if matched with
        # token.text. If found in token, skip this token when adding it to the
        # query.
        skip = {
            0: ["SP"],
        }
        s = skip.get(token.i)
        if s is not None:
            for i in s:
                if token.text.lower() == i.lower():
                    return True

        return False

    def __organisation_query(self) -> list:
        """
        Return a list of WHERE ILIKE clauses for Organisation name.

        Potential Improvements:
          - Split token.text if contains CamelCase (i.e. MelbournePool)
        """
        queries = []

        for token in self.doc:
            if self.__query_skip_token(token):
                continue

            if len(queries) == 0:
                query = token.text
            else:
                query = queries[-1] + f" {token}"
            queries.append(query)
        return queries

    def __organisation_query_postcodes(
        self, postcode: Postcode, conditions, queries
    ) -> Organisation | None:
        """
        Attempt to find an organisation by searching with conditions and
        joining on Postcode using a where id on Postcode from all Postcodes
        which match the same postcode.

        A transaction can have a postcode token which is not the same postcode in Organisation. The transaction postcode will share the same Postcode.postcode but not the same Postcode.locality as the Organisation.

        This query function will iterate over all found postcodes and attempt
        to query Organisation on similar postcodes checking for a match using conditions.
        """

        for __postcode in self.db.session.scalars(
            select(Postcode)
            .where(Postcode.postcode == postcode.postcode)
            .where(Postcode.id != postcode.id)
        ).all():
            organisations = self.db.session.scalars(
                select(Organisation, BusinessCode, Postcode)
                .join_from(Organisation, BusinessCode)
                .join_from(Organisation, Postcode)
                .where(Postcode.id == __postcode.id)
                .where(or_(*conditions))
            ).all()

            # Order is important of queries. It's important to start matching
            # on the longest query first. This is done by using sorted() with
            # key=len (length of item in list) which sorts shortest first and
            # then reversing the new list.
            for s in reversed(sorted(queries, key=len)):
                p = re.compile(rf"^{s}\S+$", re.IGNORECASE)
                for organisation in organisations:
                    if p.match(organisation.name):
                        return organisation
        return None

    def organisation(self) -> Organisation | None:
        if self.__organisation is not None:
            return self.__organisation

        postcode = self.postcode()
        ret = None
        queries = self.__organisation_query()
        stmt = select(Organisation, BusinessCode).join_from(Organisation, BusinessCode)

        if postcode is not None:
            stmt = (
                select(Organisation, BusinessCode, Postcode)
                .join_from(Organisation, BusinessCode)
                .join_from(Organisation, Postcode)
                .where(Postcode.id == postcode.id)
            )

        conditions = []
        from sqlalchemy import and_, or_

        for query in reversed(queries):
            search = f"{query}%".lower()
            conditions.append(Organisation.name.ilike(search))

        stmt = stmt.where(or_(*conditions))
        rows = self.db.session.scalars(stmt).all()
        ret = None
        for row in rows:
            if ret is not None:
                break
            for s in queries:
                if row.name.lower() == s.lower():
                    ret = row
                    break

        if ret is None and postcode is not None:
            ret = self.__organisation_query_postcodes(postcode, conditions, queries)

        if ret is not None:
            self.__organisation = ret

        return ret

    def address(self) -> str | None:
        organisation = self.__organisation
        if organisation is None:
            return None
        address = organisation.address

        if organisation.postcode:
            address = f"{address}, {organisation.postcode.locality}, {organisation.postcode.state.name}"

        return address

    def business_code(self) -> BusinessCode | None:
        organisation = self.__organisation
        if organisation is None:
            return None
        return organisation.business_code

    def anzsic(self):
        organisation = self.__organisation
        if organisation is None:
            return None
        return organisation.anzsic
