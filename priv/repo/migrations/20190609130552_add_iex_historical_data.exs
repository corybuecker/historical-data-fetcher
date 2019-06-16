defmodule HistoricalData.Repo.Migrations.AddIexHistoricalData do
  use Ecto.Migration

  def change do
    create table(:iex_historical_data, primary_key: false) do
      add(:id, :uuid, primary_key: true)
      add(:symbol, :string, null: false)
      add(:date, :date, null: false)
      add(:close, :numeric, null: false)
      add(:high, :numeric, null: false)
      add(:low, :numeric, null: false)
      add(:open, :numeric, null: false)
      add(:volume, :numeric, null: false)
      add(:unadjusted_high, :numeric, null: false)
      add(:unadjusted_low, :numeric, null: false)
      add(:unadjusted_open, :numeric, null: false)
      add(:unadjusted_close, :numeric, null: false)
      add(:unadjusted_volume, :numeric, null: false)

      timestamps()
    end

    create(unique_index(:iex_historical_data, [:symbol, :date]))
  end
end
