FROM docker.io/python:3.11-slim

ENV PATH "/root/.local/bin:$PATH"

ENV PSYCOPG_IMPL "c"

WORKDIR /app

RUN apt-get update && apt-get install -y \
    curl \
    build-essential \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/* && \
    curl -sSL https://install.python-poetry.org | python3 - --version 1.5.1

COPY poetry.lock pyproject.toml .

RUN poetry install --no-dev

RUN poetry run python -m nltk.downloader punkt

COPY boot.sh app.py . 

COPY txacct ./txacct

COPY migrations ./migrations

RUN chmod a+x boot.sh

ENTRYPOINT ["./boot.sh"]

EXPOSE 5000

