defmodule HistoricalData.Repo.Migrations.AddNasdaqCompanies do
  use Ecto.Migration

  def change do
    create table(:nasdaq_companies, primary_key: false) do
      add(:id, :uuid, primary_key: true)
      add(:symbol, :string, null: false)
      add(:name, :string, null: false)
      add(:last_sale, :numeric, null: true)
      add(:market_cap, :numeric, null: true)
      add(:adr_tso, :string, null: true)
      add(:ipo_year, :string, null: true)
      add(:sector, :string, null: true)
      add(:industry, :string, null: true)
      add(:summary_quote, :string, null: true)

      timestamps()
    end

    create(unique_index(:nasdaq_companies, [:symbol]))
  end
end
