module github.com/goverland-labs/goverland-core-feed

go 1.22

toolchain go1.23.2

replace github.com/goverland-labs/goverland-core-feed/protocol => ./protocol

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/goverland-labs/goverland-core-feed/protocol v0.0.0
	github.com/goverland-labs/goverland-platform-events v0.3.10
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/nats-io/nats.go v1.30.2
	github.com/prometheus/client_golang v1.18.0
	github.com/rs/zerolog v1.29.0
	github.com/s-larionov/process-manager v0.0.1
	github.com/shopspring/decimal v1.3.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.69.4
	google.golang.org/protobuf v1.36.3
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.1
)

require (
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
