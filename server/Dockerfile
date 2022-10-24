FROM rust:latest AS builder

WORKDIR /usr/src/server

COPY ./Cargo.lock ./Cargo.lock
COPY ./Cargo.toml ./Cargo.toml

COPY . .

RUN cargo build

FROM scratch
COPY --from=builder /usr/src/server .
USER 1000
CMD ["./server"]
