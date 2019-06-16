FROM elixir:1.8.1-alpine

COPY mix.exs mix.lock /app/

WORKDIR /app

RUN mix local.hex --force && \
  mix local.rebar --force && \
  mix deps.get && \
  mix deps.compile

COPY . /app/
RUN mix compile

CMD ["mix", "fetch"]