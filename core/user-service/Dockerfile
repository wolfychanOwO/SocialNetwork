FROM python:3.12

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

RUN mkdir /app

WORKDIR /app

COPY pyproject.toml poetry.lock ./
COPY app ./app
COPY migrations ./migrations
COPY alembic.ini ./

RUN pip install poetry && poetry config virtualenvs.create false && poetry install --no-root

EXPOSE 8001

CMD ["sh", "-c", "poetry run alembic upgrade head && poetry run uvicorn app.main:app --host 0.0.0.0 --port 8001"]

