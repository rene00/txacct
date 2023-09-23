from sqlalchemy.orm import mapped_column, Mapped, relationship
from sqlalchemy import String, ForeignKey, Integer
from typing import Any, List
from flask_sqlalchemy import SQLAlchemy

db: SQLAlchemy = SQLAlchemy()


class Transaction(db.Model):
    id: Mapped[int] = mapped_column(primary_key=True)
    memo: Mapped[str] = mapped_column(String(), nullable=False)


class Postcode(db.Model):
    __table_args__: tuple[Any] = (
        db.UniqueConstraint("postcode", "locality", name="postcode_locality"),
    )

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
