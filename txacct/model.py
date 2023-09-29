from sqlalchemy.orm import mapped_column, Mapped, relationship
from sqlalchemy import String, ForeignKey, Integer
from typing import Any, List
from flask_sqlalchemy import SQLAlchemy
import json

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
    organisations: Mapped[List["Organisation"]] = relationship(
        back_populates="postcode"
    )


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
    id: Any = db.Column(Integer(), primary_key=True)
    name: Any = db.Column(String(), nullable=False, unique=True)
    postcodes: Mapped[List["Postcode"]] = relationship(back_populates="state")


class OrganisationSource(db.Model):
    id: Any = db.Column(Integer(), primary_key=True)
    name: Mapped[str] = db.Column(String(), nullable=False, unique=True)
    organisations: Mapped[List["Organisation"]] = relationship(
        back_populates="organisation_source"
    )


class Organisation(db.Model):
    __tablename__: str = "organisation"
    __table_args__: tuple[Any] = (
        # organisation_source_id_organisation_source_id: ensure organisation from source is unique.
        db.UniqueConstraint(
            "source_id",
            "organisation_source_id",
            name="organisation_source_id_organisation_source_id",
        ),
    )
    id: Mapped[int] = mapped_column(primary_key=True)

    # source_id: the ID given to the orgnisation by the upstream data provider.
    # Not to be confused with organisation_source_id.
    source_id: Mapped[int] = mapped_column(Integer(), nullable=True)

    name: Mapped[str] = mapped_column(String(), nullable=False)
    abn: Mapped[str] = db.Column(String(), nullable=True, unique=False)
    address: Mapped[str] = db.Column(String(), nullable=True, unique=False)
    organisation_source_id: Mapped[int] = mapped_column(
        ForeignKey("organisation_source.id"), nullable=False
    )
    organisation_source: Mapped["OrganisationSource"] = relationship(
        back_populates="organisations"
    )
    anzsic_id: Mapped[int] = mapped_column(ForeignKey("anzsic.id"), nullable=True)
    anzsic: Mapped["ANZSIC"] = relationship(back_populates="organisations")
    business_code_id: Mapped[int] = mapped_column(
        ForeignKey("business_code.id"), nullable=False
    )
    business_code: Mapped["BusinessCode"] = relationship(back_populates="organisations")
    postcode_id: Mapped[int] = mapped_column(ForeignKey("postcode.id"), nullable=True)
    postcode: Mapped["Postcode"] = relationship(back_populates="organisations")


class ANZSIC(db.Model):
    __tablename__: str = "anzsic"
    __table_args__: tuple[Any] = (
        db.UniqueConstraint("code", "description", name="code_description"),
    )

    id: Mapped[int] = mapped_column(primary_key=True)

    # ANZSIC Code: the code for the anzsic item. Data from provider is showing
    # that ANZSIC items can have duplicate codes with different descriptions so
    # unique must be set to False here.
    code: Mapped[int] = db.Column(String(), nullable=True, unique=False)

    # ANZSIC Description: the description for the anzsic item. Data from
    # provider is showing that ANZSIC items can have duplicate descriptions
    # with different codes so unique must be set to False here.
    description: Mapped[int] = db.Column(String(), nullable=True, unique=False)

    organisations: Mapped[List["Organisation"]] = relationship(back_populates="anzsic")


class BusinessCode(db.Model):
    __tablename__: str = "business_code"
    __table_args__: tuple[Any] = (
        db.UniqueConstraint(
            "code", "description", name="business_code_code_description"
        ),
    )

    id: Mapped[int] = mapped_column(primary_key=True)
    code: Mapped[int] = db.Column(String(), nullable=True, unique=False)
    description: Mapped[int] = db.Column(String(), nullable=True, unique=False)
    organisations: Mapped[List["Organisation"]] = relationship(
        back_populates="business_code"
    )
