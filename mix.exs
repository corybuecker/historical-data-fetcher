defmodule HistoricalData.MixProject do
  use Mix.Project

  def project do
    [
      app: :historical_data,
      version: "1.0.0",
      elixir: "~> 1.9",
      start_permanent: Mix.env() == :prod,
      deps: deps(),
      default_release: :main,
      releases: [
        main: [
          include_executables_for: [:unix],
          applications: [runtime_tools: :permanent],
          include_system_libs: true
        ]
      ]
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
      {:ecto_sql, "~> 3.1"},
      {:postgrex, ">= 0.0.0"},
      {:floki, "~> 0.21"},
      {:timex, "~> 3.6"},
      {:nimble_csv, "~> 0.3"},
      {:credo, "~> 1.1", only: [:dev, :test], runtime: false}
    ]
  end
end
