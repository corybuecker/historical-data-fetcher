defmodule HistoricalData.NasdaqCompany do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}

  schema "nasdaq_companies" do
    field(:symbol, :string)
    field(:name, :string)
    field(:last_sale, :decimal)
    field(:market_cap, :decimal)
    field(:adr_tso, :string)
    field(:ipo_year, :string)
    field(:sector, :string)
    field(:industry, :string)
    field(:summary_quote)

    timestamps()
  end

  def changeset(nasdaq_company, params \\ %{}) do
    nasdaq_company
    |> cast(params, [
      :symbol,
      :name,
      :last_sale,
      :market_cap,
      :adr_tso,
      :ipo_year,
      :sector,
      :industry,
      :summary_quote
    ])
    |> validate_required([
      :symbol,
      :name
    ])
    |> unique_constraint(:symbol)
    |> update_change(:symbol, &String.upcase/1)
    |> update_change(:adr_tso, &replace_na/1)
    |> update_change(:ipo_year, &replace_na/1)
  end

  defp replace_na(value) do
    case value do
      "n/a" -> nil
      anything -> anything
    end
  end
end
