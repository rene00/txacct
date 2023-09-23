from .model import State, Postcode
from typing import List
from nltk.tokenize import word_tokenize
from sqlalchemy import func, select
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

    def postcode(self, locality: str, state: State) -> List[Postcode]:
        return (
            self.db.session.query(Postcode)
            .filter(Postcode.locality == locality, Postcode.state_id == state.id)
            .all()
        )
