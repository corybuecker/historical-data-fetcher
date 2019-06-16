defmodule HistoricalData.Repo.Migrations.AddIexSymbols do
  use Ecto.Migration

  def change do
    create table(:iex_symbols, primary_key: false) do
      add(:id, :uuid, primary_key: true)
      add(:currency, :string, null: false)
      add(:date, :date, null: false)
      add(:iex_id, :string, null: false)
      add(:is_enabled, :boolean, null: false)
      add(:name, :string, null: false)
      add(:region, :string, null: false)
      add(:symbol, :string, null: false)
      add(:type, :string, null: false)
      add(:exchange, :string)

      timestamps()
    end

    create(unique_index(:iex_symbols, :iex_id))
  end
end
