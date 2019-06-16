defmodule Polygon do
  use HTTPoison.Base

  def process_request_url(url) do
    with %{path: path, query: query} <- URI.parse(url),
         query <- if(query, do: query, else: ""),
         query <-
           URI.decode_query(query)
           |> Map.merge(%{"apiKey" => Application.get_env(:historical_data, :polygon_key)}) do
      URI.parse("https://api.polygon.io")
      |> URI.merge(path)
      |> URI.merge("?" <> URI.encode_query(query))
      |> URI.to_string()
    end
  end

  def process_response_body(body) do
    body
    |> Poison.decode!()
  end
end
