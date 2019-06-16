defmodule Mix.Tasks.UpdateIexHistoricalData do
  use Mix.Task

  require Logger

  alias HistoricalData.IexHistoricalDataSymbol
  alias HistoricalData.Repo
  alias HistoricalData.IexUpsertSymbolDate

  @shortdoc "Simply runs the HistoricalData.run/0 function"
  def run(_) do
    Mix.Task.run("app.start")

    Ecto.Adapters.SQL.query!(
      Repo,
      "REFRESH MATERIALIZED VIEW iex_historical_data_symbols;"
    )

    Timex.Interval.new(
      from: Timex.shift(Timex.today(), days: -2),
      until: Timex.today(),
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
    |> Enum.each(fn interval -> IexUpsertSymbolDate.upsert(company.symbol, interval) end)

    :timer.sleep(500)
  end
end
