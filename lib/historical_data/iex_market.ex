defmodule HistoricalData.IexMarket do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}

  schema "iex_markets" do
    field(:name, :string, primary_key: true)
    field(:mic, :string)
    field(:tape_id, :string)
    field(:oats_id, :string)
    field(:type, :string)
    field(:long_name, :string)

    timestamps()
  end

  def changeset(iex_market, params \\ %{}) do
    iex_market
    |> cast(params, [:name, :mic, :tape_id, :oats_id, :type, :long_name])
    |> validate_required([:name, :mic, :type])
    |> unique_constraint(:name)
  end
end
