defmodule HistoricalData.MixProject do
  use Mix.Project

  def project do
    [
      app: :historical_data,
      version: "0.1.0",
      elixir: "~> 1.8",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [mod: {HistoricalData, []}, extra_applications: [:logger]]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:httpoison, "~> 1.5"},
      {:poison, "~> 4.0"},
      {:ecto_sql, "~> 3.0"},
      {:postgrex, ">= 0.0.0"},
      {:floki, "~> 0.21"},
      {:timex, "~> 3.1"},
      {:nimble_csv, "~> 0.3"},
      {:credo, "~> 1.0", only: [:dev, :test], runtime: false}
    ]
  end
end
