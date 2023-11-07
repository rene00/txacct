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

    tm = TransactionMeta(memo=transaction.memo, db=db)

    response: Dict[str, Any] = {
        "id": transaction.id,
        "memo": transaction.memo,
    }

    response["locality"] = {}

    state = tm.state()
    if state is not None:
        response["locality"]["state"] = state.name

    postcode = tm.postcode()
    if postcode is not None:
        response["locality"]["postcode"] = postcode.postcode
        response["locality"]["name"] = postcode.locality

    organisation = tm.organisation()
    if organisation is not None:
        response["organisation"] = organisation.name

    if not bool(response["locality"]):
        del response["locality"]

    address = tm.address()
    if address is not None:
        response["address"] = address

    business_code = tm.business_code()
    if business_code is not None:
        response["description"] = business_code.description

    return jsonify(response)
