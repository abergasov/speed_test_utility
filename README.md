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


# Assignment

## Description
Create a small GO library that tests the download and upload speeds by using Ookla's https://www.speedtest.net/ and Netflix's https://fast.com/.

## Requirements

1. The library has 1 exposed API, that starts the speed test by providing the choice between the 2 speed providers(Ookla's & Netflix's) and returns the Mbps for both download and upload
2. The test coverage for the implementation has to be at least 80%
3. Provide 2 benchmark tests (1 for each of the 2, Ookla's & Netflix's) for the exposed API
4. Provide documentation for your implementation, both via code docs and README.md

## Evaluation

We look for:
1. The way you structure your code
2. The code quality (code standards, best practices, design patterns, etc)
3. The way you write tests
4. The way you write docs
5. Attention to details

## Suggestions

Please don't reinvent the wheel. The internet has everything you need. We care about how you solve puzzles and how you use the legos to build the end product.

You can use whichever 3rd party libraries you see fit to get the job done.

## Deadline

It has to be submitted within 24hours since the moment you received this assignment.

## Deliverable

Create a public personal github repo where you commit and push this assignment, and provide the link to it.