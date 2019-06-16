defmodule HistoricalData.Repo.Migrations.AddSp500Companies do
  use Ecto.Migration

  def change do
    create table(:sp500_companies, primary_key: false) do
      add(:id, :uuid, primary_key: true)
      add(:symbol, :string, null: false)
      add(:name, :string, null: false)

      timestamps()
    end

    create(unique_index(:sp500_companies, [:symbol]))
  end
end
