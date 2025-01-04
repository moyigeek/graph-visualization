# Visualize_dep_graph
Visualize Package Dependency Graph

## Usage

### install
install go package
```shell
cd  graph_server
go mod tidy
```

set config with your own database
```shell
cp config.toml.example config.toml
```

### run

```shell
cd  graph_server
go run main.go
```

```shell    
cd  graph_client
python -m http.server <port>
```

### open browser
open browser and input url
```shell
http://localhost:<port>
```

