# speed_test_utility

## speedtest
```shell
internal/service/speedtest/speedtest.go
```
uses `github.com/showwin/speedtest-go/speedtest` under the hood

forked to `github.com/abergasov/speedtest-go` and add ability pass http client as interface. It allow write tests|benchmarks much easier

simply call object download|upload methods

## fast.com
```shell
internal/service/fast/fast.go
```
uses `github.com/gesquive/fast-cli/fast` under the hood

library via headless chromium get token to make requests

than simply run download|upload http requests

### upload
service generate mock payload of fixed size. that run this to upload endpoint

measure spending time and calculate total upload bytes

### download
make several concurrent requests to download endpoint and wait when response body send `EOF`

measure spending time and calculate total consumed bytes

## usage 
run testing network speed via speedtest.net
```shell
make speed
```
via fast.com
```shell
make fast
```

#### Dev
run tests:
```shell
make test
```

run benchamrks:
```shell
make bench
```

run linters
```shell
make lint
```