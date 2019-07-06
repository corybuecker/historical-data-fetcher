FROM elixir:1.9.0-alpine as compiler
COPY . /app/
WORKDIR /app
RUN mix local.hex --force && \
  mix local.rebar --force && \
  mix deps.get && \
  mix release

FROM alpine:3.9
RUN apk add ncurses-libs
COPY --from=compiler /app/_build/dev/rel/historical_data /app
WORKDIR /app
ENTRYPOINT ["/app/bin/historical_data", "eval"]