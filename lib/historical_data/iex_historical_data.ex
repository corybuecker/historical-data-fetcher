defmodule HistoricalData.IexHistoricalData do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}

  schema "iex_historical_data" do
    field(:symbol, :string)
    field(:date, :date)
    field(:close, :decimal)
    field(:high, :decimal)
    field(:low, :decimal)
    field(:open, :decimal)
    field(:volume, :decimal)
    field(:unadjusted_high, :decimal)
    field(:unadjusted_low, :decimal)
    field(:unadjusted_open, :decimal)
    field(:unadjusted_close, :decimal)
    field(:unadjusted_volume, :decimal)

    timestamps()
  end

  def changeset(iex_historical_data, params \\ %{}) do
    iex_historical_data
    |> cast(params, [
      :symbol,
      :date,
      :close,
      :high,
      :low,
      :open,
      :volume,
      :unadjusted_close,
      :unadjusted_low,
      :unadjusted_open,
      :unadjusted_high,
      :unadjusted_volume
    ])
    |> validate_required([
      :symbol,
      :date,
      :close,
      :high,
      :low,
      :open,
      :volume,
      :unadjusted_close,
      :unadjusted_low,
      :unadjusted_open,
      :unadjusted_high,
      :unadjusted_volume
    ])
    |> unique_constraint(:symbol, name: :iex_historical_data_symbol_date_index)
    |> update_change(:symbol, &String.upcase/1)
  end
end
