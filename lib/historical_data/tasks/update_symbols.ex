defmodule HistoricalData.Tasks.UpdateSymbols do
  require Logger

  alias HistoricalData.IexSymbol
  alias HistoricalData.Repo

  def run() do
    Application.ensure_all_started(:historical_data)

    {:ok, %{body: symbols}} = Iexcloud.get("/stable/ref-data/symbols")
    symbols |> Enum.each(&insert_symbol/1)
  end

  defp insert_symbol(symbol) do
    Logger.info(symbol |> inspect())

    {:ok, _} =
      IexSymbol.changeset(%IexSymbol{}, %{
        currency: Map.get(symbol, "currency"),
        date: Map.get(symbol, "date"),
        exchange: Map.get(symbol, "exchange"),
        is_enabled: Map.get(symbol, "isEnabled"),
        iex_id: Map.get(symbol, "iexId"),
        name: Map.get(symbol, "name"),
        region: Map.get(symbol, "region"),
        symbol: Map.get(symbol, "symbol"),
        type: Map.get(symbol, "type")
      })
      |> Repo.insert(on_conflict: :replace_all_except_primary_key, conflict_target: :iex_id)
  end
end
