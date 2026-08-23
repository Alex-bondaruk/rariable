# Rariable service

Go client and HTTP service for two Rarible Protocol endpoints: `GetOwnershipByID` and `QueryTraitsWithRarity`.

## 1. Setup

Get a free API key at https://api.rarible.org/dashboard.

```bash
cp .env.example .env
# fill in RARIBLE_API_KEY
```

## 2. Test

```bash
go test ./...
```

## 3. Run

```bash
set -a; . ./.env; set +a
go run ./cmd/server
```

```bash
curl localhost:8080/healthz

curl "localhost:8080/ownerships/ETHEREUM:0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d:664:0x4459084da2d3a774c436f2e75f2e3fe9335dc5de"

curl -X POST localhost:8080/traits/rarity \
  -H "Content-Type: application/json" \
  -d '{"collectionId":"ETHEREUM:0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d","properties":[{"key":"Hat","value":"Halo"}]}'
```

## 4. Docker

```bash
docker build -t rariable:0.1.0 .
docker run -p 8080:8080 --env-file .env rariable:0.1.0
```

## 5. Helm

```bash
helm lint deploy/helm/rariable
helm template rariable deploy/helm/rariable --set secret.apiKey=your-key
```
