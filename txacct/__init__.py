from flask import Flask
from flask_migrate import Migrate
from .blueprints.transactions import transactions_bp
from .cli import postcode_cli, nltk_cli
from .model import db


def create_app(test_config: dict | None = None) -> Flask:
    app: Flask = Flask(__name__)
    if test_config:
        app.config.from_mapping(test_config)
    else:
        app.config.from_prefixed_env(prefix="TXACCT")
    db.init_app(app)
    Migrate(app, db)
    app.register_blueprint(transactions_bp)
    app.cli.add_command(postcode_cli)
    app.cli.add_command(nltk_cli)
    return app
