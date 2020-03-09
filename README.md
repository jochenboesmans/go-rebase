[![jochenboesmans](https://circleci.com/gh/jochenboesmans/go-rebase.svg?style=svg)](https://app.circleci.com/pipelines/github/jochenboesmans/go-rebase)

# Example usage

Here's an example of a simple fx market being rebased in USD.

## Input

```json
{
  "rebaseId": "USD",
  "maxPathLength": 3,
  "market": [
    {
      "baseId": "USD",
      "quoteId": "EUR",
      "exchangeMarkets": [
        {
          "currentBid": 1.12,
          "currentAsk": 1.13,
          "baseVolume": 100
        }
      ]
    },
    {
      "baseId": "EUR",
      "quoteId": "GBP",
      "exchangeMarkets": [
        {
          "currentBid": 1.15,
          "currentAsk": 1.16,
          "baseVolume": 100
        }
      ]
    },
    {
      "baseId": "GBP",
      "quoteId": "JPY",
      "exchangeMarkets": [
        {
          "currentBid": 0.0073,
          "currentAsk": 0.0074,
          "baseVolume": 100
        }
      ]
    }
  ]
}
```

## Output

```json
{
  "rebaseId": "USD",
  "market": [
    {
      "baseId": "USD",
      "quoteId": "EUR",
      "exchangeMarkets": {
        "FX1": {
          "currentBid": 1.12,
          "currentAsk": 1.13,
          "baseVolume": 100
        }
      }
    },
    {
      "baseId": "EUR",
      "quoteId": "GBP",
      "exchangeMarkets": [
        {
          "currentBid": 1.29375,
          "currentAsk": 1.305,
          "baseVolume": 112.5
        }
      ]
    },
    {
      "baseId": "GBP",
      "quoteId": "JPY",
      "exchangeMarkets": [
        {
          "currentBid": 0.0094854375,
          "currentAsk": 0.009615375,
          "baseVolume": 129.9375
        }
      ]
    }
  ]
}
```
