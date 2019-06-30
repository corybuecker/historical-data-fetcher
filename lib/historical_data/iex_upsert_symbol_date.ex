defmodule HistoricalData.IexUpsertSymbolDate do
  @base_path "/stable/stock/%s/chart/date/%s?chartByDay=true"

  require Logger

  alias HistoricalData.IexHistoricalData
  alias HistoricalData.Repo
  alias Ecto.Changeset

  def upsert(symbol, date) do
    {symbol, date} |> skip_if_in_database() |> fetch() |> insert()
  end

  defp skip_if_in_database(q) do
    case IexHistoricalData.exists?(q) do
      true -> {:error, "skipping #{q |> Tuple.to_list() |> Enum.join(", ")} since it exists"}
      false -> {:ok, q}
    end
  end

  defp fetch({:ok, {symbol, date}}) do
    case @base_path
         |> add_placeholder(symbol)
         |> add_placeholder(Date.to_iso8601(date, :basic))
         |> Iexcloud.get() do
      {:ok, %{body: historical_data}} -> {:ok, {symbol, date, historical_data}}
      err -> {:error, err}
    end
  end

  defp fetch({:error, error}) do
    {:error, error}
  end

  defp insert({:ok, {symbol, _date, historical_data}}) do
    with do
      historical_data |> Enum.each(fn datum -> insert_historical_data(symbol, datum) end)
    else
      err -> err |> IO.inspect()
    end
  end

  defp insert({:error, error}) do
    IO.inspect(error)
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
