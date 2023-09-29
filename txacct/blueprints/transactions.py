from flask import Blueprint, jsonify, Response, request, abort
from txacct.model import db
from txacct.model import Transaction, Postcode
from txacct.transaction_meta import TransactionMeta
from typing import Dict, Any, List

transactions_bp: Blueprint = Blueprint(
    "transactions", __name__, url_prefix="/transactions"
)


@transactions_bp.route("/", methods=["POST"])
def transactions() -> Response:
    if not request.is_json:
        abort(400)

    data = request.get_json()

    memo = data.get("memo", None)
    if memo is None:
        abort(400)

    transaction: Transaction = Transaction(memo=memo)
    db.session.add(transaction)
    db.session.commit()

    tm = TransactionMeta(memo=transaction.memo, db=db)

    response: Dict[str, Any] = {
        "id": transaction.id,
        "memo": transaction.memo,
    }

    locality: Dict[str, Dict] = {"names": {}}
    state = tm.state()
    if state is not None:
        locality["state"] = dict(
            {
                "name": state.name,
            }
        )

        postcodes: List[Postcode] = tm.postcode()
        if len(postcodes) >= 1:
            names = dict()
            for i in postcodes:
                postcode = i[0]
                locality_name = {
                    "name": postcode.locality,
                    "postcode": postcode.postcode,
                }

                if postcode.sa3:
                    locality_name["sa3"] = {"name": postcode.sa3.name}

                if postcode.sa4:
                    locality_name["sa4"] = {"name": postcode.sa4.name}

                weight: int = 0
                for i in range(100, 1, -1):
                    if i in names:
                        continue
                    weight = i
                    break

                names[weight] = locality_name

            locality["names"] = names

    response["locality"] = locality

    organisation = tm.organisation()
    if organisation is not None:
        response["organisation"] = organisation.name

    address = tm.address()
    if address is not None:
        response["address"] = address

    business_code = tm.business_code()
    if business_code is not None:
        response["description"] = business_code.description

    return jsonify(response)
