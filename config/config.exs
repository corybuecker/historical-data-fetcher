use Mix.Config

config :historical_data, HistoricalData.Repo, database: "historical_data"

config :historical_data,
  polygon_key: System.get_env("APCA_API_KEY_ID"),
  iexcloud_key: System.get_env("IEXCLOUD_KEY")

config :historical_data, ecto_repos: [HistoricalData.Repo]
