defmodule Mix.Tasks.UpdateMarkets do
  use Mix.Task

  require Logger

  alias HistoricalData.IexMarket
  alias HistoricalData.Repo

  @shortdoc "Simply runs the HistoricalData.run/0 function"
  def run(_) do
    Mix.Task.run("app.start")
    {:ok, %{body: markets}} = Iexcloud.get("/stable/ref-data/market/us/exchanges")
    markets |> Enum.each(&insert_market/1)
  end

  defp insert_market(market) do
    Logger.info(market |> inspect())

    {:ok, _} =
      IexMarket.changeset(%IexMarket{}, %{
        name: Map.get(market, "name"),
        mic: Map.get(market, "mic"),
        tape_id: Map.get(market, "tapeId"),
        oats_id: Map.get(market, "oatsId"),
        type: Map.get(market, "type"),
        long_name: Map.get(market, "longName")
      })
      |> Repo.insert(on_conflict: :replace_all_except_primary_key, conflict_target: :name)
  end
end
