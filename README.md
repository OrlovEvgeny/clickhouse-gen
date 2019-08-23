# ClickHouse model generator

<p align="center"><img src="https://raw.githubusercontent.com/OrlovEvgeny/clickhouse-gen/master/build/logo.png"></p>
<p align="center">

**Install clickhouse-gen:**
````bash
curl -s https://raw.githubusercontent.com/OrlovEvgeny/clickhouse-gen/master/install.sh | bash -s --
````

**check that everything is fine**
````bash
~ $ clickhouse-gen -help

    Usage of clickhouse-gen:
      -c string
            path to config file. example --c=config.yaml (default "./config.yaml")
      -pack string
            go package. example --pack=model (default "model")
      -path string
            path to output folder. example --path=./model (default "./")
      -table string
            target table. example --table=simple_table (default "simple")

````

**Example use**
````bash
~ $ clickhouse-gen --path=./src/entity --pack=entity --table=example_table --c=config.yaml
````

**result:**

[Generated](https://raw.githubusercontent.com/OrlovEvgeny/clickhouse-gen/master/example/example.table.model.go)

# Build Source

**For compilation you need to install >= [Golang1.8](https://medium.com/@patdhlk/how-to-install-go-1-8-on-ubuntu-16-04-710967aa53c9)**

```bash
~ $ make build
```

# License:
[MIT](LICENSE)