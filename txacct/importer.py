from txacct.model import StatisticalArea3, StatisticalArea4, State, Postcode
from sqlalchemy.exc import IntegrityError


class Importer:
    def __init__(self, session, **kwargs) -> None:
        self.session = session
        self.postcode_data: dict = kwargs.get("postcode_data", {})

        if self.postcode_data is None:
            raise Exception("postcode_data must be set")

    def do(self) -> None:
        """import locality data into db"""
        for i in self.postcode_data:
            sa3 = None
            if i.get("sa3") != "" and i.get("sa3name") != "":
                sa3 = StatisticalArea3(
                    code=i.get("sa3"),
                    name=i.get("sa3name"),
                )
                self.session.add(sa3)
                try:
                    self.session.commit()
                except IntegrityError:
                    self.session.rollback()
                    sa3 = (
                        self.session.query(StatisticalArea3)
                        .filter_by(code=int(sa3.code), name=i.get("sa3name"))
                        .one()
                    )

            sa4 = None
            if i.get("sa4") != "" and i.get("sa4name") != "":
                sa4 = StatisticalArea4(
                    code=i.get("sa4"),
                    name=i.get("sa4name"),
                )
                self.session.add(sa4)
                try:
                    self.session.commit()
                except IntegrityError:
                    self.session.rollback()
                    sa4 = (
                        self.session.query(StatisticalArea4)
                        .filter_by(code=int(sa4.code), name=sa4.name)
                        .one()
                    )

            state = None
            if i.get("state") != "":
                state = State(name=i.get("state"))
                self.session.add(state)
                try:
                    self.session.commit()
                except IntegrityError:
                    self.session.rollback()
                    state = self.session.query(State).filter_by(name=state.name).one()

            postcode: Postcode = Postcode(
                postcode=i.get("postcode"),
                locality=i.get("locality"),
                state=state,
                sa3=sa3,
                sa4=sa4,
            )

            self.session.add(postcode)
            try:
                self.session.commit()
            except IntegrityError:
                self.session.rollback()
