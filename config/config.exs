use Mix.Config

config :historical_data, HistoricalData.Repo, database: "historical_data"
config :historical_data, ecto_repos: [HistoricalData.Repo]
