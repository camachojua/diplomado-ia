import yfinance as yf

# Fortune 500 tickers (example: top 5)
tickers = ["AAPL", "MSFT", "AMZN", "GOOGL", "TSLA"]  # Replace with actual Fortune 500 tickers

# Download stock price data for the last 4 years
data = yf.download(tickers, start="2020-01-01", end="2023-12-31", group_by="ticker")

# Display the first few rows
for ticker in tickers:
    print(f"\n{ticker} Stock Data:")
    print(data[ticker].head())
