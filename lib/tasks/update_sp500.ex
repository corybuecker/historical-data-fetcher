defmodule Mix.Tasks.UpdateSp500 do
  use Mix.Task

  require Logger

  alias HistoricalData.Sp500Company
  alias HistoricalData.Repo

  @shortdoc "Simply runs the HistoricalData.run/0 function"
  def run(_) do
    Mix.Task.run("app.start")
    Repo.query("TRUNCATE sp500_companies", [])
    fetch_html() |> parse_html() |> extract_companies() |> Enum.each(&insert_symbol/1)
  end

  defp insert_symbol(symbol) do
    Logger.info(symbol |> inspect())

    {:ok, _} =
      Sp500Company.changeset(%Sp500Company{}, %{
        symbol: Enum.at(symbol, 0),
        name: Enum.at(symbol, 1)
      })
      |> Repo.insert(on_conflict: :replace_all_except_primary_key, conflict_target: :symbol)
  end

  defp extract_companies(elements) do
    elements
    |> Enum.map(fn row -> row |> extract_company() end)
    |> Enum.filter(fn [a, b] -> a != nil && b != nil end)
  end

  defp extract_company({_, _, contents}) do
    contents
    |> Enum.take(2)
    |> extract_nested()
    |> extract_nested()
  end

  defp extract_nested(ary) do
    ary
    |> Enum.map(fn a ->
      case a do
        {_, _, [val]} -> val
        _ -> nil
      end
    end)
  end

  defp parse_html(html) do
    Floki.find(html, "table#constituents tr")
  end

  defp fetch_html do
    case Wikipedia.get("/wiki/List_of_S%26P_500_companies") do
      {:ok, %{body: html}} -> html
      _ -> ""
    end
  end
end
