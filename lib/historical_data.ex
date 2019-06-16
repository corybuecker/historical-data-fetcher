defmodule HistoricalData do
  require Logger
  use Application

  def start(_type, _args) do
    Polygon.start()
    Iexcloud.start()
    Wikipedia.start()
    Nasdaq.start()

    children = [
      HistoricalData.Repo
    ]

    Supervisor.start_link(children, strategy: :one_for_one)
  end
end
