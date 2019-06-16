defmodule HistoricalData.Sp500Company do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}

  schema "sp500_companies" do
    field(:symbol, :string)
    field(:name, :string)

    timestamps()
  end

  def changeset(sp500_company, params \\ %{}) do
    sp500_company
    |> cast(params, [
      :symbol,
      :name
    ])
    |> validate_required([
      :symbol,
      :name
    ])
    |> unique_constraint(:symbol)
    |> update_change(:symbol, &String.upcase/1)
  end
end
