defmodule HistoricalData.Tasks.UpdateIexHistoricalData do
  require Logger

  alias HistoricalData.IexHistoricalDataSymbol
  alias HistoricalData.Repo
  alias HistoricalData.IexUpsertSymbolDate

  def run(args \\ []) do
    Application.ensure_all_started(:historical_data)

    Ecto.Adapters.SQL.query!(
      Repo,
      "REFRESH MATERIALIZED VIEW iex_historical_data_symbols;"
    )

    from =
      case args do
        [] -> Timex.shift(Timex.today(), days: -2)
        [f] -> Timex.shift(Timex.today(), days: f)
        [f, _t] -> Timex.shift(Timex.today(), days: f)
      end

    to =
      case args do
        [] -> Timex.today()
        [_f] -> Timex.today()
        [_f, t] -> Timex.shift(Timex.today(), days: t)
      end

    Logger.info("Starting from #{from} until #{to}")

    Timex.Interval.new(
      from: from,
      until: to,
      left_open: false,
      right_open: false
    )
    |> process()
  end

  defp process(intervals) do
    IexHistoricalDataSymbol
    |> Repo.all()
    |> Enum.each(fn company -> upsert_historical_data(company, intervals) end)
  end

  defp upsert_historical_data(company, intervals) do
    Logger.info(company |> inspect())
    Logger.info(intervals |> inspect())

    intervals
    |> Enum.each(fn interval ->
      IexUpsertSymbolDate.upsert(company.symbol, interval)
      :timer.sleep(50)
    end)
  end
end
