defmodule HistoricalData.IexSymbol do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}

  schema "iex_symbols" do
    field(:currency, :string)
    field(:date, :date)
    field(:iex_id, :string)
    field(:is_enabled, :boolean)
    field(:name, :string)
    field(:region, :string)
    field(:symbol, :string)
    field(:type, :string)
    field(:exchange, :string)

    timestamps()
  end

  def changeset(iex_symbol, params \\ %{}) do
    iex_symbol
    |> cast(params, [
      :currency,
      :date,
      :iex_id,
      :is_enabled,
      :name,
      :region,
      :symbol,
      :type,
      :exchange
    ])
    |> validate_required([
      :currency,
      :date,
      :iex_id,
      :is_enabled,
      :name,
      :region,
      :symbol,
      :type
    ])
    |> unique_constraint(:iex_id)
  end
end
