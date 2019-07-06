defmodule HistoricalData.Tasks.UpdatePreviousIexHistoricalData do
  require Logger

  alias HistoricalData.IexHistoricalDataSymbol
  alias HistoricalData.Repo
  alias HistoricalData.IexUpsertSymbolDate

  def run() do
    Application.ensure_all_started(:historical_data)

    Ecto.Adapters.SQL.query!(
      Repo,
      "REFRESH MATERIALIZED VIEW iex_historical_data_symbols;"
    )

    IexHistoricalDataSymbol
    |> Repo.all()
    |> Enum.each(fn company -> upsert_historical_data(company) end)
  end

  defp upsert_historical_data(company) do
    Logger.info(company |> inspect())

    IexUpsertSymbolDate.upsert_previous(company.symbol)
    :timer.sleep(50)
  end
end
