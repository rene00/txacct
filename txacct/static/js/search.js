// static/js/people.js

import { sendForm } from "./request.js";

export class Search {
  constructor() {
    this.activeSearchForm();
  }

  activeSearchForm() {
    const searchForm = document.querySelector(".search-card form");
    new CreateSearchForm(searchForm);
  }
}

class CreateSearchForm {
  constructor(el) {
    this.form = el;
    this.searchButton = el.querySelector("button[data-action='search']");
    this.searchButton.addEventListener(
      "click",
      this.handleSearchClick.bind(this)
    );
  }

  handleSearchClick(event) {
    event.preventDefault();
    sendForm(this.form, "POST", "/transactions/", this.handleTransactionResponse);
  }

  handleTransactionResponse(rawData) {
    const data = JSON.parse(rawData);

    const transactionCard = document.querySelector(".transaction-card");
    const transactionContent = transactionCard.querySelector(".transaction-content");

    const notFound = transactionContent.querySelector(".not-found");

    const organisation = transactionContent.querySelector(".organisation");
    if (data.hasOwnProperty("organisation")) {
      organisation.textContent = "Organisation: " + data.organisation;
      organisation.classList.remove("hidden");

      notFound.classList.add("hidden");
    } else {
      organisation.classList.add("hidden");

      notFound.textContent = "Not found";
      notFound.classList.remove("hidden");
    }

    const address = transactionContent.querySelector(".address");
    if (data.hasOwnProperty("address")) {
      address.textContent = "Address: " + data.address;
      address.classList.remove("hidden");
    } else {
      address.classList.add("hidden");
    }

    const description = transactionContent.querySelector(".description");
    if (data.hasOwnProperty("description")) {
      description.textContent = "Description: " + data.description;
      description.classList.remove("hidden");
    } else {
      description.classList.add("hidden");
    }
  }
}
