# Oracle Price Feeder Dockerfile

## Build Docker Image
Change `VERSION` to the release you want to build.

```bash
VERSION=main
git clone https://github.com/eni-chain/go-eni.git
cd oracle/price-feeder/dockerfile || exit
docker build --build-arg VERSION=$VERSION -t price-feeder:latest .
```

## Create `config.toml`
Edit your `address`, `validator`, `grpc_endpoint`, `tmrpc_endpoint` you may need to modifify your firewall to allow this container to reach your chain-node. See [offical docs](https://docs.kujira.app/validators/run-a-node/oracle-price-feeder) for more details.

```bash
sudo tee config.toml <<EOF
gas_adjustment = 1.5
gas_prices = "0.00125ueni"
enable_server = true
enable_voter = true
provider_timeout = "500ms"

[server]
listen_addr = "0.0.0.0:7171"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[deviation_thresholds]]
base = "USDT"
threshold = "2"

[account]
address = "eni..."
chain_id = "go-eni"
validator = "enivaloper..."
prefix = "eni"

[keyring]
backend = "file"
dir = "/root/.eni"

[rpc]
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"
tmrpc_endpoint = "http://localhost:26657"

[telemetry]
enable_hostname = true
enable_hostname_label = true
enable_service_label = true
enabled = true
global_labels = [["chain_id", "kaiyo-1"]]
service_name = "price-feeder"
type = "prometheus"
prometheus_retention = 120

[[provider_endpoints]]
name = "binance"
rest = "https://api1.binance.com"
websocket = "stream.binance.com:9443"

[[currency_pairs]]
base = "ATOM"
chain_denom = "uatom"
providers = [
  "binance",
  "kraken",
  "coinbase",
]
quote = "USD"
EOF
```

## Create `client.toml`
change node to your favorite `rpc` node

```bash
sudo tee client.toml <<EOF
chain-id = "go-eni"
keyring-backend = "file"
output = "text"
node = "tcp://localhost:26657"
broadcast-mode = "sync"
EOF
```

## Recover oracle `keyring-file` to local file
```bash
enid keys add oracle --keyring-backend file --recover
```
In the eni home directory (~/.eni/) you should see the `keyring-file` folder.  This will be mounted as a volume when running the docker container.

## Run Docker Image
```bash
docker run \
--env PRICE_FEEDER_PASS=password \
-v ~/.eni/keyring-file:/root/.eni/keyring-file \
-v "$PWD"/config.toml:/root/price-feeder/config.toml \
-v "$PWD"/client.toml:/root/.eni/config/client.toml \
-it price-feeder /root/price-feeder/config.toml
```
