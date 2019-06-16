defmodule HistoricalData.IexUpsertSymbolDate do
  @base_path "/stable/stock/%s/chart/date/%s?chartByDay=true"

  require Logger

  alias HistoricalData.IexHistoricalData
  alias HistoricalData.Repo
  alias Ecto.Changeset

  def upsert(symbol, date) do
    with {:ok, %{body: historical_data}} <-
           @base_path
           |> add_placeholder(symbol)
           |> add_placeholder(Date.to_iso8601(date, :basic))
           |> Iexcloud.get() do
      historical_data |> Enum.each(fn d -> insert_historical_data(symbol, d) end)
    else
      err -> Logger.error(err)
    end
  end

  defp insert_historical_data(symbol, historical_data) do
    Logger.info(historical_data |> inspect())

    {:ok, _} =
      IexHistoricalData.changeset(%IexHistoricalData{}, %{
        symbol: symbol,
        date: Map.get(historical_data, "date"),
        close: Map.get(historical_data, "close"),
        high: Map.get(historical_data, "high"),
        low: Map.get(historical_data, "low"),
        open: Map.get(historical_data, "open"),
        volume: Map.get(historical_data, "volume"),
        unadjusted_close: Map.get(historical_data, "uClose"),
        unadjusted_low: Map.get(historical_data, "uLow"),
        unadjusted_open: Map.get(historical_data, "uOpen"),
        unadjusted_high: Map.get(historical_data, "uHigh"),
        unadjusted_volume: Map.get(historical_data, "uVolume")
      })
      |> Repo.insert(
        on_conflict: :replace_all_except_primary_key,
        conflict_target: [:symbol, :date]
      )
  end

  defp add_placeholder(string, value) do
    string |> String.replace("%s", value, global: false)
  end
end
