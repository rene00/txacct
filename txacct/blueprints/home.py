from flask import Blueprint, render_template

home_bp: Blueprint = Blueprint("home", __name__, url_prefix="/")


@home_bp.route("/", methods=["GET", "POST"])
def home() -> str:
    return render_template("home.html", current_page="home")


@home_bp.route("/about", methods=["GET"])
def about() -> str:
    return render_template("about.html", current_page="about")
