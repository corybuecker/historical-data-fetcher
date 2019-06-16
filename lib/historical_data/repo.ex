defmodule HistoricalData.Repo do
  use Ecto.Repo,
    otp_app: :historical_data,
    adapter: Ecto.Adapters.Postgres
end
