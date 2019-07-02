defmodule Mix.Tasks.UpdateNasdaqCompanies do
  use Mix.Task

  require Logger

  alias HistoricalData.NasdaqCompany
  alias HistoricalData.Repo
  alias NimbleCSV.RFC4180, as: CSV

  def run(_) do
    Mix.Task.run("app.start")
    Repo.query("TRUNCATE nasdaq_companies", [])

    fetch_csv("/screening/companies-by-industry.aspx?industry=Technology&render=download")
    |> parse_csv()
    |> Enum.each(&insert_symbol/1)

    fetch_csv("/screening/companies-by-industry.aspx?industry=Public+Utilities&render=download")
    |> parse_csv()
    |> Enum.each(&insert_symbol/1)
  end

  defp insert_symbol(symbol) do
    Logger.info(symbol |> inspect())

    with {:ok, _} <-
           NasdaqCompany.changeset(%NasdaqCompany{}, %{
             symbol: Enum.at(symbol, 0),
             name: Enum.at(symbol, 1),
             last_sale: Enum.at(symbol, 2),
             market_cap: Enum.at(symbol, 3),
             adr_tso: Enum.at(symbol, 4),
             ipo_year: Enum.at(symbol, 5),
             sector: Enum.at(symbol, 6),
             industry: Enum.at(symbol, 7),
             summary_quote: Enum.at(symbol, 8)
           })
           |> Repo.insert(on_conflict: :replace_all_except_primary_key, conflict_target: :symbol) do
      true
    else
      err -> err |> IO.inspect()
    end
  end

  defp parse_csv(html) do
    CSV.parse_string(html)
  end

  defp fetch_csv(path) do
    case Nasdaq.get(path) do
      {:ok, %{body: html}} -> html
      _ -> ""
    end
  end
end
