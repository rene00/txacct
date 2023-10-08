from .model import ANZSIC, Organisation, State, Postcode, BusinessCode
from typing import List
from nltk.tokenize import word_tokenize
from sqlalchemy import func, select, or_, and_
from sqlalchemy.exc import NoResultFound


class TransactionMeta:
    def __init__(self, memo: str, db) -> None:
        self.memo: str = memo
        self.db = db
        self.tokenized: List[str] = word_tokenize(self.memo)

    def state(self) -> State | None:
        """state returns the Australian state of the transaction."""
        l: str = self.tokenized[-1]
        for row in self.db.session.execute(select(State)).all():
            state = row[0]
            if l.lower().endswith(state.name.lower()):
                return state

        return None

    def postcode(self) -> List[Postcode] | None:
        conditions = []
        seen = []
        memo = list(reversed(self.tokenized))

        state = self.state()
        if state is None:
            return None

        # flag which is set when a state name (e.g. VIC) is stripped from
        # tokenized memo.
        flag_state = False

        for i, val in enumerate(memo):
            # skip state.
            if val == state.name:
                continue

            # if val ends in a state, strip the state from val.
            if val.endswith(state.name) and i < 2:
                val = val.rstrip(state.name)
                flag_state = True

            # skip first item.
            if i == (len(memo) - 1):
                continue

            # if 2nd item in reverse list, append last iterated item. From
            # ["SOUTH", "WHARF"], this should return "SOUTH WHARF". This will
            # limit or conditions to length of 2.
            if i == 2 or (i == 1 and flag_state):
                previous = memo[i - 1].rstrip(state.name)
                val = f"{val} {previous}"
                flag_state = False

            o = Postcode.locality == val
            conditions.append(o)
            seen.append(val)

        return self.db.session.execute(
            select(Postcode).where(
                and_(or_(*conditions), Postcode.state_id == state.id)
            )
        ).all()

    def organisation(self) -> Organisation | None:
        # The name to use when searching for the organisation.
        name = self.memo

        state = self.state()
        if state is not None:
            # Remove the state from the organisation name.
            name = name.replace(state.name, "")

        stmt = (
            select(Organisation, BusinessCode)
            .where(func.lower(Organisation.name) == name.rstrip().lower())
            .join_from(Organisation, BusinessCode)
        )

        postcode = self.postcode()
        if postcode is not None and len(postcode) >= 1:
            # Remove the locality from the organisation name.
            name = name.replace(postcode[0][0].locality, "").rstrip()

            # Include the postcode on the join when selecting for organisation
            # to increase accuracy.
            stmt = (
                select(Organisation, BusinessCode, Postcode)
                .join_from(Organisation, BusinessCode)
                .join_from(Organisation, Postcode)
                .where(func.lower(Organisation.name) == name.lower())
                .where(Postcode.id == postcode[0][0].id)
            )

        return self.db.session.scalars(stmt).first()

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
