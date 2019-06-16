defmodule HistoricalData.Repo.Migrations.AddIexMarkets do
  use Ecto.Migration

  def change do
    create table(:iex_markets, primary_key: false) do
      add(:id, :uuid, primary_key: true)
      add(:name, :string, null: false)
      add(:mic, :string, null: false)
      add(:tape_id, :string)
      add(:oats_id, :string)
      add(:type, :string, null: false)
      add(:long_name, :string)

      timestamps()
    end

    create(unique_index(:iex_markets, :name))
  end
end
