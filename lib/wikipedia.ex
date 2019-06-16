defmodule Wikipedia do
  use HTTPoison.Base

  def process_request_url(url) do
    with %{path: path, query: query} <- URI.parse(url),
         query <- if(query, do: query, else: ""),
         query <- URI.decode_query(query) do
      URI.parse("https://en.wikipedia.org/wiki")
      |> URI.merge(path)
      |> URI.merge("?" <> URI.encode_query(query))
      |> URI.to_string()
    end
  end

  def process_response_body(body) do
    body
  end
end
