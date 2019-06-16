defmodule HistoricalData.Repo.Migrations.AddIexSymbolsListView do
  use Ecto.Migration

  def up do
    execute("""
    CREATE MATERIALIZED VIEW iex_historical_data_symbols AS
      SELECT DISTINCT
        *
      FROM (
        SELECT
          symbol
        FROM
          nasdaq_companies
        UNION
        SELECT
          symbol
        FROM
          sp500_companies
      ) as inq1;
    ;
    """)
  end

  def down do
    execute("""
      DROP MATERIALIZED VIEW iex_historical_data_symbols;
    """)
  end
end
