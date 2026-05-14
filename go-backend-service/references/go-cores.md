# go-cores Deep Notes

Source: `/Users/nuttawitbutsara/projects/nb/go-cores`

Module:
```go
module github.com/witwoywhy/go-cores
```

---

## Overview

`go-cores` is the shared runtime toolkit for Go services. It provides config loading, app metadata, structured logging, tracing, HTTP server helpers, error mapping, outbound HTTP clients, Kafka pub/sub, GORM initialization/logging, circuit breaker setup, storage clients, crypto helpers, and small utilities.

Most packages follow this shape:

1. Config is loaded by `vipers.Init()`.
2. Runtime config is read from Viper with `viper.UnmarshalKey(...)`.
3. Service code initializes infrastructure in `main.go`.
4. Handlers/services receive `logger.Logger`.
5. Logs carry `trace_id` and `span_id`.

---

---

## Package Relationship Summary

```mermaid
flowchart LR
    vipers --> apps
    vipers --> logs
    vipers --> gins
    vipers --> reqs
    vipers --> kafka
    vipers --> gorms
    vipers --> storage
    vipers --> circuitbreaker

    apps --> logs
    apps --> gins
    apps --> reqs
    apps --> kafka

    logger --> logs
    logger --> gins
    logger --> reqs
    logger --> kafka
    logger --> gorms

    errs --> gins
    httpserve --> gins
    contexts --> gins
    tracers --> logs
    tracers --> gins
    pubsub --> kafka
    pubsub --> logs
    dbs --> gorms
    cryptos --> kafka
    utils --> logs
    utils --> gins
    utils --> reqs
```

---

## Config / Bootstrap

Packages that load service configuration and app metadata.

---

## `vipers` — Configuration Loader

### Purpose

Loads `config.yaml` into global Viper state, enables env override, and normalizes all keys back into Viper.

### Files

```text
vipers/config.go
vipers/vipers.go
```

### Config Source

Defaults:

```text
VIPER_PATH_CONFIG=./config
VIPER_FILENAME_CONFIG=config
VIPER_FILE_EXTENSION_CONFIG=yaml
```

Equivalent default file:

```text
./config/config.yaml
```

### API

```go
vipers.SetPathConfig("./config")
vipers.SetFileNameConfig("config")
vipers.SetFileExtensionConfig("yaml")
vipers.Init()
```

### Behavior

```go
viper.SetConfigName(fileNameConfig)
viper.SetConfigType(fileExtensionConfig)
viper.AddConfigPath(pathConfig)
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
viper.AutomaticEnv()
viper.ReadInConfig()
```

Env override example:

```text
app.name       -> APP_NAME
log.tracer_url -> LOG_TRACER_URL
db.transfer.dsn -> DB_TRANSFER_DSN
```

### Example

```go
func init() {
    vipers.Init()
}
```

### Notes

- Must run before packages that read config with `viper.UnmarshalKey`.
- Panics if config file cannot be read.

---

---

## `apps` — App Metadata and Shared Log Keys

### Purpose

Loads app identity and defines shared constants used by logs, tracing, HTTP middleware, and message correlation.

### Files

```text
apps/config.go
apps/const.go
apps/helper.go
```

### Config

```yaml
app:
  name: example
  env: dev
  time_zone: UTC
```

### API

```go
apps.Init()
fmt.Println(apps.Config.Name)
```

Config struct:

```go
type ConfigInfo struct {
    Name     string `mapstructure:"name"`
    Env      string `mapstructure:"env"`
    TimeZone string `mapstructure:"time_zone"`
}
```

### Important Constants

```go
apps.TraceID
apps.SpanID
apps.Authorization
apps.Message
apps.StartInbound
apps.EndInbound
apps.SummaryInbound
apps.StartOutbound
apps.EndOutbound
apps.SummaryOutbound
```

### Header Masking

```go
apps.MaskHeader(header)
```

Masks headers present in `apps.HeaderMaskingList`:

```go
authorization
x-api-key
```

### Notes

- `apps.Init()` is normally base infrastructure for all services.
- `apps` constants are used across `logs`, `gins`, `reqs`, and Kafka tracing.

---

---

## Observability

Packages for logging, log interfaces, and distributed tracing.

---

## `logger` — Logging Interface

### Purpose

Defines the interface consumed by service/adaptor/repository code.

### File

```text
logger/logger.go
```

### Interface

```go
type Logger interface {
    Info(obj any)
    Debug(obj any)
    Warn(obj any)
    Error(err any)

    Infof(format string, obj ...any)
    Debugf(format string, obj ...any)
    Warnf(format string, obj ...any)
    Errorf(format string, obj ...any)

    JSON(m map[string]any)
}
```

### Usage

All application layers can accept this interface instead of importing `logs.Log` directly:

```go
func Execute(request *Request, l logger.Logger) (*Response, errs.Error) {
    l.Infof("processing request")
    return &Response{}, nil
}
```

---

---

## `logs` — Structured Logs, Span Logs, Producer Logs

### Purpose

Provides structured JSON logging with trace/span fields, optional OpenTelemetry tracing, masking, and optional log publishing to a pubsub producer.

### Files

```text
logs/config.go
logs/init.go
logs/log.go
logs/option.go
logs/json-handler.go
logs/producer-handler.go
logs/helper.go
logs/const.go
```

### Config

```yaml
log:
  level: info              # info | debug | warn | error
  masking_list: password|token|secret
  tracer_url: http://localhost:4318
```

All possible config fields:

```go
type ConfigInfo struct {
    Level       string `mapstructure:"level"`
    MaskingList string `mapstructure:"masking_list"`
    TracerUrl   string `mapstructure:"tracer_url"`

    IsEnableTracer bool
}
```

`IsEnableTracer` is runtime state, not config input.

### Init

Stdout JSON logging:

```go
logs.Init()
```

Publish logs to Kafka/pubsub:

```go
logs.Init(
    logs.AddProducer(kafka.NewProducer(kafka.AddConfigKey("log.producer"))),
)
```

Shutdown:

```go
defer logs.Shutdown()
```

### Logger Creation

```go
l := logs.New(map[string]any{
    apps.TraceID: "trace-id",
    apps.SpanID:  "span-id",
})

l.Info("message")
l.Infof("hello %s", "world")
l.JSON(map[string]any{"field": "value"})
```

Root log:

```go
l, span := logs.NewLog()
if span != nil {
    defer span.End()
}
```

Span action:

```go
l, end := logs.NewSpanLogAction(l, "REPO CREATE TRANSACTION")
defer end()
```

### Behavior

- If `log.tracer_url` is set, `tracers.Init(...)` runs.
- `NewSpanLog` and `NewSpanLogAction` switch to tracer-backed implementations.
- Logs include service name from `app.name`.
- `masking_list` is split by `|` and recursively masks matching JSON fields.
- `ProducerHandler` sends log fields to configured pubsub producer.

### Relation Diagram

```mermaid
flowchart LR
    logs.Init --> Config[log config]
    Config --> Stdout[JsonHandler stdout]
    Config --> Producer[ProducerHandler pubsub]
    Config --> Tracer[tracers.Init if tracer_url]
    logs.NewLog --> Logger[logger.Logger]
    Logger --> Span[NewSpanLogAction]
    Span --> AppCode[service/adaptor/repository]
```

---

---

## `tracers` — OpenTelemetry Setup

### Purpose

Initializes global OpenTelemetry tracer provider for logs and HTTP request spans.

### File

```text
tracers/init.go
```

### Config

Indirectly configured through:

```yaml
log:
  tracer_url: http://localhost:4318
```

### API

```go
tracers.Init("http://localhost:4318")
defer tracers.Shutdown()
```

### Behavior

- Uses OTLP HTTP exporter.
- Service name comes from `app.name`.
- Environment comes from `app.env`.
- `Shutdown()` is protected by `sync.Once`.

---

---

## HTTP

Packages for HTTP server config, Gin app wrapping, request context, and error responses.

---

## `http-serve` — HTTP Server Config Struct

### Purpose

Defines the config shape consumed by `gins.New()`.

### File

```text
http-serve/config.go
```

### Config

```yaml
http_serve:
  port: 8080
  ignore_log_body:
    - /v1/health
  error_code_mapping: |
    {}
```

### Struct

```go
type HTTPServe struct {
    Port             string   `mapstructure:"port"`
    IgnoreLogBody    []string `mapstructure:"ignore_log_body"`
    ErrorCodeMapping string   `mapstructure:"error_code_mapping"`
}
```

---

---

## `gins` — HTTP App, Middleware, Handler Wrappers

### Purpose

Wraps Gin with app registration, middleware, request logging, error mapping, route-context handler wrappers, and graceful shutdown.

### Files

```text
gins/apps.go
gins/handler.go
gins/app-log-middleware.go
gins/app-error-middleware.go
gins/helper.go
```

### Config

```yaml
app:
  name: transfer-service
  env: dev

http_serve:
  port: 8080
  ignore_log_body:
    - /v1/health
  error_code_mapping: |
    {
      "40000": {
        "th": { "message": "กรุณากรอก {{ field }} ให้ครบถ้วน" },
        "en": { "message": "{{ field }} is required." }
      }
    }
```

### API

```go
app := gins.New()

app.UseMiddleware(app.Log())
app.UseMiddleware(app.Error())

app.Register(
    http.MethodGet,
    "/logger",
    app.WithLogger(func(ctx *gin.Context, l logger.Logger) {
        l.Info("HANDLE WITH LOGGER")
    }),
)

app.Register(
    http.MethodGet,
    "/context",
    app.WithRouteContext(func(ctx *gin.Context, rctx *contexts.RouteContext, l logger.Logger) {
        l.Info("HANDLE WITH ROUTE CONTEXT")
    }),
)

app.ListenAndServe(func() {
    logs.Shutdown()
})
```

### Handler Types

```go
type HandleWithLogger func(ctx *gin.Context, l logger.Logger)
type HandleWithRouteContextLogger func(ctx *gin.Context, rctx *contexts.RouteContext, l logger.Logger)
```

### Middleware Flow

```mermaid
flowchart TD
    Request --> LogMiddleware[app.Log]
    LogMiddleware --> MakeLogger[NewRequestLog + ctx.Set logger]
    MakeLogger --> Handler[WithLogger / WithRouteContext]
    Handler --> ErrorMiddleware[app.Error]
    ErrorMiddleware --> Mapping[map errs.Errs to localized message]
    ErrorMiddleware --> Response
```

### `app.Log()`

- Wraps response writer to capture response body.
- Creates request logger and span.
- Logs START INBOUND, END INBOUND, SUMMARY INBOUND.
- Masks Authorization header.
- Skips body logging for paths in `ignore_log_body`.

### `app.Error()`

- Reads `ctx.Errors`.
- Expects `*errs.Errs`.
- Applies language mapping from `http_serve.error_code_mapping`.
- Language source: `X-Language`, context value, then default `th`.
- Returns JSON with status from error.

### `WithLogger`

Use when handler/service does not need `RouteContext`:

```go
app.WithLogger(func(ctx *gin.Context, l logger.Logger) {})
```

### `WithRouteContext`

Use when handler/service needs `contexts.RouteContext`:

```go
app.WithRouteContext(func(ctx *gin.Context, rctx *contexts.RouteContext, l logger.Logger) {})
```

Requires `rctx` to exist in Gin context. If missing, it returns internal error.

---

---

## `contexts` — Header and Route Context

### Purpose

Defines common request context/header DTOs for route-aware handlers.

### Files

```text
contexts/header.go
contexts/routeContext.go
```

### Types

```go
type Header struct {
    Authorization string `header:"Authorization"`

    RequestID string `header:"X-Request-Id"`
    TraceID   string `header:"Trace-Id"`
    SpanID    string `header:"Span-Id"`

    Language language.Language `header:"Accept-Language"`
}

type RouteContext struct {
    UserRefId string
    Header
}
```

### Usage

```go
func Handle(ctx *gin.Context, rctx *contexts.RouteContext, l logger.Logger) {
    l.Infof("user: %s", rctx.UserRefId)
}
```

---

---

## `errs` — Standard Error Shape and Mapping

### Purpose

Defines application error shape, constructors, multi-error wrapper, and language-based error message mapping.

### Files

```text
errs/const.go
errs/error.go
errs/errors.go
errs/mapping.go
errs/provider.go
```

### Error Codes

```go
Err40000 = "40000"
Err40001 = "40001"
Err40004 = "40004"

Err50000 = "50000"
Err50001 = "50001"
Err50002 = "50002"
Err50003 = "50003"
```

### Error Shape

```go
type Err struct {
    Status      int    `json:"-"`
    ErrorCode   string `json:"code"`
    Message     string `json:"message"`
    Description string `json:"description"`

    Data map[string]any `json:"-"`
    Err  error          `json:"-"`
}

type Errs struct {
    Errors []Err `json:"errors"`
}
```

### Constructors

```go
errs.New(status, code)
errs.NewCustom(status, code, message, description, err)
errs.NewInternalError(err) // 500 / Err50001
errs.NewExternalError(err) // 503 / Err50003
errs.NewBadRequestError(err) // 400 / Err40000
errs.NewBusinessError("030001", err) // 409 / domain code
```

### Error Mapping Config

Used by `gins.Error()` middleware:

```yaml
http_serve:
  error_code_mapping: |
    {
      "40000": {
        "th": {
          "message": "กรุณากรอก {{ field }} ให้ครบถ้วน"
        },
        "en": {
          "message": "{{ field }} is required."
        }
      }
    }
```

### Usage

```go
if err := ctx.BindJSON(&request); err != nil {
    ctx.Error(errs.NewBadRequestError(err))
    ctx.Abort()
    return
}

if notFound {
    return nil, errs.NewBusinessError(domain.Err030004)
}
```

### Notes

- `gins.Error()` expects errors pushed into Gin to be `*errs.Errs`.
- Mapping uses `raymond` templates and `Err.Data`.

---

---

## Outbound Communication

Packages for outbound HTTP calls and circuit-breaker protection.

---

## `reqs` — Outbound HTTP Client

### Purpose

Wraps `github.com/witwoywhy/req` with config-driven client creation, request builder methods, and outbound request/response logging.

### Files

```text
reqs/config.go
reqs/client.go
reqs/request.go
reqs/response.go
```

### Config

```yaml
integrations:
  get_user:
    base_url: http://localhost:3005
    method: GET
    url: /users/{id}
    timeout: 10s
    enable_insecure_skip_verify: true
    enable_ignore_log_body: false
    api_key: API-KEY
    username:
    password:
```

All possible config fields:

```go
type Config struct {
    BaseUrl                  string        `mapstructure:"base_url"`
    Method                   string        `mapstructure:"method"`
    Url                      string        `mapstructure:"url"`
    Timeout                  time.Duration `mapstructure:"timeout"`
    EnableInsecureSkipVerify bool          `mapstructure:"enable_insecure_skip_verify"`
    EnableIgnoreLogBody      bool          `mapstructure:"enable_ignore_log_body"`
    ApiKey                   string        `mapstructure:"api_key"`
    Username                 string        `mapstructure:"username"`
    Password                 string        `mapstructure:"password"`
}
```

### API

```go
client := reqs.NewClient("integrations.get_user")

var result User
var responseError Error

resp := client.Request(l).
    SetHeader("X-Api-Key", client.Config().ApiKey).
    SetPathParam("id", "001").
    SetResult(&result).
    SetError(&responseError).
    Do()

if resp.IsErrorState() {
    l.Errorf("failed to get user: %v", responseError)
} else if resp.Error() != nil {
    l.Errorf("unknown error: %v", resp.Error())
}
```

### Request Builder

```go
SetContext(ctx)
SetPathParam(key, value)
SetPathParams(map[string]string)
AddQueryParam(key, value)
AddQueryParams(key, values...)
SetHeader(key, value)
SetHeaders(map[string]string)
SetBearerAuthToken(token)
SetBasicAuth(username, password)
SetFormData(map[string]interface{})
SetFileReader(paramName, filename, reader)
SetBody(body)
SetError(errDto)
SetResult(successDto)
Do()
```

### Logging Behavior

- Logs START OUTBOUND before request.
- Logs END OUTBOUND and SUMMARY OUTBOUND after response.
- Masks Authorization header.
- Skips body logging when `enable_ignore_log_body: true`.
- Supports GET, POST, PATCH, PUT, DELETE.

---

---

## `circuit-breaker` — Hystrix Wrapper

### Purpose

Configures and executes Hystrix circuit breakers.

### Files

```text
circuit-breaker/init.go
circuit-breaker/do.go
```

### Config

```yaml
circuit_breakers:
  cbs_account_inquiry:
    timeout: 11000
    max_concurrent_requests: 100
    request_volume_threshold: 20
    sleep_window: 15000
    error_percent_threshold: 75
```

All fields:

```text
timeout
max_concurrent_requests
request_volume_threshold
sleep_window
error_percent_threshold
```

### API

```go
circuitbreaker.Init()

err := circuitbreaker.Do("cbs_account_inquiry", func() error {
    // call external system
    return nil
}, nil)
```

Async:

```go
ch := circuitbreaker.Go("job", fn, fallback)
err := <-ch
```

### Notes

- `Do` lowercases command name.
- Handle `hystrix.CircuitError` in adaptors when you need custom timeout/open-circuit mapping.

---

---

## Messaging

Packages for generic pub/sub interfaces and Kafka implementations.

---

## `pubsub` — Interfaces

### Purpose

Defines generic producer and consumer group interfaces. `kafka` implements these.

### Files

```text
pubsub/producer.go
pubsub/consumer-group.go
pubsub/handle.go
```

### Interfaces

```go
type Producer interface {
    Produce(key string, v any, l logger.Logger) error
    Shutdown(l logger.Logger) error
}

type ConsumerGroup interface {
    Consume(l logger.Logger, fn HandlerFunc, closeFunc func())
}

type HandlerFunc = func(ctx context.Context, topic, group, key string, value []byte) bool
```

`HandlerFunc` returns `true` to ack/commit.

---

---

## `kafka` — Franz-go Producer and Consumer Group

### Purpose

Provides Kafka producer and consumer group implementations of `pubsub.Producer` and `pubsub.ConsumerGroup`.

### Files

```text
kafka/config.go
kafka/option.go
kafka/producer.go
kafka/consumer-group.go
kafka/helper.go
kafka/log.go
```

### Producer Config

```yaml
pubsub:
  pub:
    broker: localhost:9092
    topic: test
    batch_size: 1048576
    max_buffer_records: 10000
    linger: 5ms
    request_timeout: 60s
    delivery_timeout: 180s
    cert:
      type: value # value or file
      ca: |
        -----BEGIN CERTIFICATE-----
        -----END CERTIFICATE-----
      cert: |
        -----BEGIN CERTIFICATE-----
        -----END CERTIFICATE-----
      key: |
        -----BEGIN PRIVATE KEY-----
        -----END PRIVATE KEY-----
```

All producer config fields:

```go
broker
topic
cert.type
cert.ca
cert.key
cert.cert
batch_size
max_buffer_records
linger
request_timeout
delivery_timeout
```

Defaults:

```go
batch_size: 1048576
max_buffer_records: 10000
linger: 10ms
request_timeout: 60s
delivery_timeout: 180s
```

### Consumer Config

```yaml
pubsub:
  sub:
    broker: localhost:9092
    topic: test
    consumer_group: test-group
    fetch_max_bytes: 5242880
    fetch_min_bytes: 1024
    fetch_max_partition_bytes: 1048576
    fetch_max_wait: 1s
    session_timeout: 60s
    heartbeat_interval: 10s
    max_poll_interval: 600s
    request_timeout: 60s
    consume_reset_offset: end # start or end
    cert:
      type: value # value or file
      ca: |
        -----BEGIN CERTIFICATE-----
        -----END CERTIFICATE-----
      cert: |
        -----BEGIN CERTIFICATE-----
        -----END CERTIFICATE-----
      key: |
        -----BEGIN PRIVATE KEY-----
        -----END PRIVATE KEY-----
```

All consumer config fields:

```go
broker
topic
consumer_group
cert.type
cert.ca
cert.key
cert.cert
fetch_max_bytes
fetch_min_bytes
fetch_max_partition_bytes
fetch_max_wait
session_timeout
heartbeat_interval
max_poll_interval
request_timeout
consume_reset_offset
```

Default reset offset: `end`.

### Producer API

```go
p := kafka.NewProducer(
    kafka.AddConfigKey("pubsub.pub"),
    kafka.AddLogger(logs.L),
)

err := p.Produce("message-key", map[string]any{"hello": "world"}, logs.L)
defer p.Shutdown(logs.L)
```

### Consumer API

```go
c := kafka.NewConsumerGroup(
    kafka.AddConfigKey("pubsub.sub"),
    kafka.AddLogger(logs.L),
)

c.Consume(logs.L, func(ctx context.Context, topic, group, key string, value []byte) bool {
    fmt.Println(topic, group, key, string(value))
    return true
}, func() {
    logs.Shutdown()
})
```

### Consumer Flow

```mermaid
flowchart LR
    Poll[PollFetches] --> Records[EachRecord]
    Records --> Handler[HandlerFunc]
    Handler -->|true| Commit[MarkCommitRecords]
    Handler -->|false| Skip[No commit]
    Signal[SIGINT/SIGTERM] --> Close[Client.Close + closeFunc]
```

### Notes

- TLS cert config supports `type: value` or file paths.
- Producer marshals non-string/non-byte payloads as JSON.
- Producer uses LZ4 compression and all ISR acks.

---

---

## Database

Packages for database config and GORM integration.

---

## `dbs` — Database Config

### Purpose

Defines DB config used by `gorms.Init`.

### Files

```text
dbs/const.go
dbs/db.go
```

### Config

```yaml
db:
  transfer:
    driver: pg
    dsn: postgres://root:root@localhost:5432/witwoywhy
    host: localhost
    port: "5432"
    username: root
    password: root
    database: witwoywhy
    timeout: 1m
    read_timeout: 1m
    write_timeout: 1m
    max_idle_conns: 1
    max_conns: 5
    max_life_time: 1h
```

All fields:

```go
driver
dsn
host
port
username
password
database
timeout
read_timeout
write_timeout
max_idle_conns
max_conns
max_life_time
```

Supported drivers:

```go
mysql
pg
postgres
```

### Notes

The example config in `gorms/example/configs/config.yaml` uses camelCase keys like `maxIdleConns`, while the struct tags use snake_case like `max_idle_conns`. Prefer the struct-supported snake_case keys.

`ToDSN()` behavior was fixed in `go-cores v1.0.2`: explicit `dsn` is used when provided; otherwise it builds a DSN from host/port/user/password/database fields.

```go
if d.DSN != "" {
    return d.DSN
}
// otherwise build DSN from fields
```

---

---

## `gorms` — GORM Initialization and SQL Logging

### Purpose

Initializes GORM DB connections from config and provides GORM logger integration with `logger.Logger`.

### Files

```text
gorms/init.go
gorms/log.go
gorms/log-option.go
gorms/helper.go
gorms/const.go
```

### API

```go
db := gorms.Init("db.transfer")
```

Create GORM session with structured logger:

```go
l := logs.New(map[string]any{})
spl := logs.NewSpanLog(l)
gormLog := gorms.NewGormLog(gorms.AddLogger(spl))

tx := db.Session(&gorm.Session{Logger: gormLog})
```

Debug SQL to stdout:

```go
gormLog := gorms.NewGormLog(
    gorms.AddLogger(l),
    gorms.Debug(),
)
```

### Behavior

- Supports PostgreSQL and MySQL.
- Pings DB on init.
- Applies max idle/open/lifetime connection settings.
- `NewGormLog` logs DB operations as outbound actions.
- SQL table name is extracted from query with regex.

### Notes

- Non-debug logger logs start/end and table name, not full SQL.
- Several GORM logger interface methods panic if called directly (`Info`, `Warn`, `Error`, `LogMode`); current use relies on `Trace`.

---

---

## Storage

Package for unified object storage clients.

---

## `storage` — S3/R2/MinIO, GCS, Azure Blob

### Purpose

Provides a unified object storage interface.

### Files

```text
storage/config.go
storage/iface.go
storage/request-info.go
storage/s3.go
storage/gcs.go
storage/azure.go
```

### Interface

```go
type Storager interface {
    Upload(ctx context.Context, request *UploadRequest) error
    GetUrl(filePath string) string
    GetSignedUrl(ctx context.Context, filePath string, ttl time.Duration) (string, error)
    Delete(ctx context.Context, filePath string) error
}
```

Upload request:

```go
type UploadRequest struct {
    FilePath         string
    File             io.Reader
    ContentType      string
    DownloadFileName *string
}
```

### Config

All possible fields:

```yaml
storage:
  files:
    target: s3 # s3 | r2 | minio | blob | gcs
    bucketName: my-bucket
    prefix: uploads
    publicUrl: https://cdn.example.com
    region: ap-southeast-1

    # S3 / R2 / MinIO
    accessKey: access
    secretKey: secret
    endpoint: https://s3.example.com
    acl: public-read

    # Azure Blob
    azureConnection: DefaultEndpointsProtocol=...

    # GCS
    serviceAccount: |
      {"type":"service_account"}
```

### Init

```go
store := storage.Init("storage.files")
```

### Usage

```go
err := store.Upload(ctx, &storage.UploadRequest{
    FilePath:    "reports/a.pdf",
    File:        reader,
    ContentType: "application/pdf",
})

url := store.GetUrl("reports/a.pdf")
signed, err := store.GetSignedUrl(ctx, "reports/a.pdf", 15*time.Minute)
err = store.Delete(ctx, "reports/a.pdf")
```

### Target Behavior

- `s3`, `r2`, `minio` use AWS S3 client. Non-`s3` targets set `BaseEndpoint`.
- `gcs` uses service account JSON and optional prefix.
- `blob` uses Azure connection string and prefix as container name.

---

---

## Utilities

Small shared utility packages and constants.

---

## `cryptos` — AES, Base64, TLS Config

### Purpose

Small crypto helpers for AES-CBC compatible with CryptoJS, base64 encode/decode, and TLS config creation for Kafka.

### Files

```text
cryptos/aes.go
cryptos/b64.go
cryptos/tls.go
```

### AES CryptoJS-Compatible

```go
key := "12345678901234561234567890123456" // 32 bytes for AES-256
iv := []byte("1234567812345678")          // 16 bytes

aes := cryptos.NewAesCompatibleCryptoJS(key, iv)
enc, err := aes.Encrypt(map[string]any{"hello": "world"})
dec, err := aes.Decrypt(enc)
```

Supported encrypt input:

```text
map[string]any
string
int / int8 / int16 / int32 / int64
float32 / float64
```

### Base64

```go
encoded := cryptos.B64Encode([]byte("hello"))
decoded, err := cryptos.B64Decode(encoded)
```

### TLS

From files:

```go
tlsConfig, err := cryptos.NewTLSConfigFromFile("cert.pem", "key.pem", "ca.pem")
```

From string values:

```go
tlsConfig, err := cryptos.NewTLSConfig(certPEM, keyPEM, caPEM)
```

Used by Kafka producer/consumer TLS setup.

---

---

## `enum` — Shared Constants

### Purpose

Provides shared enum values.

### Files

```text
enum/language/language.go
enum/platform/platform.go
```

### Language

```go
language.TH // "th"
language.EN // "en"
```

Used by error mapping and language detection.

### Platform

```go
platform.Android
platform.IOS
platform.Web
```

---

---

## `utils` — Small Generic Helpers

### Purpose

Utility helpers used across packages.

### Files

```text
utils/make-map.go
utils/masking.go
utils/not-nil.go
utils/slice-array.go
```

### `MakeMap`

Builds a map from a slice of structs using field names.

```go
type Data struct {
    Id   int64
    Name string
}

items := []Data{{Id: 1, Name: "A"}, {Id: 2, Name: "B"}}

byID := utils.MakeMap[int64, string](items, "Id", "Name")
exists := utils.MakeMap[int64, bool](items, "Id", "true")
structs := utils.MakeMap[int64, Data](items, "Id", "struct")
```

### `MaskString`

```go
utils.MaskString("secret", "*") // s****t
```

### `MaskHeader`

```go
utils.MaskHeader(apps.Authorization, header)
```

Sets header value to `***` if present.

### `NotNil`

```go
value := utils.NotNil(ptr)
```

Returns zero value when pointer is nil.

### `SliceArray`

```go
chunks := utils.SliceArray(items, 100)
```

Splits a slice into chunks of `maxPerSlice`.

---

---

## Full Config Reference

```yaml
app:
  name: service-name
  env: dev
  time_zone: UTC

log:
  level: info
  masking_list: password|token|secret
  tracer_url: http://localhost:4318
  producer:
    broker: localhost:9092
    topic: log
    batch_size: 1048576
    max_buffer_records: 10000
    linger: 10ms
    request_timeout: 60s
    delivery_timeout: 180s
    cert:
      type: file
      ca: /path/ca.pem
      cert: /path/cert.pem
      key: /path/key.pem

http_serve:
  port: 8080
  ignore_log_body:
    - /v1/health
  error_code_mapping: |
    {
      "40000": {
        "th": { "message": "กรุณากรอก {{ field }} ให้ครบถ้วน" },
        "en": { "message": "{{ field }} is required." }
      }
    }

db:
  transfer:
    driver: pg
    dsn: postgres://user:pass@localhost:5432/db
    host: localhost
    port: "5432"
    username: user
    password: pass
    database: db
    timeout: 1m
    read_timeout: 1m
    write_timeout: 1m
    max_idle_conns: 5
    max_conns: 10
    max_life_time: 1h

integrations:
  cbs_account_inquiry:
    base_url: http://localhost:3005
    method: GET
    url: /account/{accountId}
    timeout: 10s
    enable_insecure_skip_verify: true
    enable_ignore_log_body: false
    api_key: key
    username:
    password:

circuit_breakers:
  cbs_account_inquiry:
    timeout: 11000
    max_concurrent_requests: 100
    request_volume_threshold: 20
    sleep_window: 15000
    error_percent_threshold: 75

pubsub:
  pub:
    broker: localhost:9092
    topic: events
    batch_size: 1048576
    max_buffer_records: 10000
    linger: 10ms
    request_timeout: 60s
    delivery_timeout: 180s
    cert:
      type: file
      ca: /path/ca.pem
      cert: /path/cert.pem
      key: /path/key.pem
  sub:
    broker: localhost:9092
    topic: events
    consumer_group: service-group
    fetch_max_bytes: 5242880
    fetch_min_bytes: 1024
    fetch_max_partition_bytes: 1048576
    fetch_max_wait: 1s
    session_timeout: 60s
    heartbeat_interval: 10s
    max_poll_interval: 600s
    request_timeout: 60s
    consume_reset_offset: end
    cert:
      type: file
      ca: /path/ca.pem
      cert: /path/cert.pem
      key: /path/key.pem

storage:
  files:
    target: s3
    bucketName: my-bucket
    prefix: uploads
    publicUrl: https://cdn.example.com
    region: ap-southeast-1
    accessKey: access
    secretKey: secret
    endpoint: https://s3.example.com
    acl: public-read
    azureConnection:
    serviceAccount:
```
