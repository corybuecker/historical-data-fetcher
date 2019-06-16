defmodule Mix.Tasks.UpdateSymbols do
  use Mix.Task

  require Logger

  alias HistoricalData.IexSymbol
  alias HistoricalData.Repo

  @shortdoc "Simply runs the HistoricalData.run/0 function"
  def run(_) do
    Mix.Task.run("app.start")
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
