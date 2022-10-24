FROM rust:latest

WORKDIR /usr/src/server

COPY ./Cargo.lock ./Cargo.lock
COPY ./Cargo.toml ./Cargo.toml

COPY . .

RUN cargo build

EXPOSE 8000

CMD ["cargo", "run"]
