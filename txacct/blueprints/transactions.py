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

        postcode: List[Postcode] = tm.postcode(locality=tm.tokenized[-2], state=state)
        if len(postcode) >= 1:
            names = dict()
            for i in postcode:
                locality_name = {"name": i.locality, "postcode": i.postcode}

                if i.sa3:
                    locality_name["sa3"] = {"name": i.sa3.name}

                if i.sa4:
                    locality_name["sa4"] = {"name": i.sa4.name}

                weight: int = 0
                for i in range(100, 1, -1):
                    if i in names:
                        continue
                    weight = i
                    break

                names[weight] = locality_name

            locality["names"] = names

    response["locality"] = locality

    return jsonify(response)
