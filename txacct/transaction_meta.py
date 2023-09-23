from .model import State, Postcode
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
