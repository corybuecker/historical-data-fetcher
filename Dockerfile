FROM elixir:1.9.0-alpine as compiler
COPY mix.exs mix.lock /app/
WORKDIR /app
RUN mix local.hex --force && \
  mix local.rebar --force && \
  mix deps.get && \
  mix deps.compile
COPY . /app/
RUN mix compile && \
  mix release

FROM alpine:3.9
RUN apk add ncurses-libs
COPY --from=compiler /app/_build/dev/rel/main /app/main
WORKDIR /app
ENTRYPOINT ["/app/main/bin/main", "eval"]