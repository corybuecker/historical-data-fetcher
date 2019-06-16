defmodule HistoricalData.IexHistoricalDataSymbol do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key false

  schema "iex_historical_data_symbols" do
    field(:symbol, :string)
  end
end
