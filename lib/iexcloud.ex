defmodule Iexcloud do
  use HTTPoison.Base
  require Logger

  def process_request_url(url) do
    with %{path: path, query: query} <- URI.parse(url),
         query <- if(query, do: query, else: ""),
         query <-
           URI.decode_query(query)
           |> Map.merge(%{
             "token" => Application.get_env(:historical_data, :iexcloud_key)
           }) do
      URI.parse("https://cloud.iexapis.com")
      |> URI.merge(path)
      |> URI.merge("?" <> URI.encode_query(query))
      |> URI.to_string()
    end
  end

  def process_response_body(body) do
    with {:ok, results} <- body |> Poison.decode() do
      results
    else
      err -> handle_error(err)
    end
  end

  defp handle_error(err) do
    err |> IO.inspect()

    []
  end
end
