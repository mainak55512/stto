# STTO

Command-line utility written in Go to check total line of code in a file present in a directory.



## Authors

- [Mainak Bhattacharjee](https://github.com/mainak55512)


## Dependencies

- go 1.22.5
- github.com/mattn/go-runewidth v0.0.9
- github.com/olekukonko/tablewriter v0.0.5

## Benchmark

#### Benchmark was run on the clone of '[Redis](https://github.com/redis/redis)' repository

![Demo](./resources/benchmark.gif)
**N.B: stto is no way near the more established options like 'scc' or 'tokei' in terms of features. It is in early development stage and isn't production ready.

All the tools read over 1.5k files
![stto](./resources/stto_redis.png)
![scc](./resources/scc_redis.png)
![tokei](./resources/tokei_redis.png)

## Usage
### usage 1:
![stto_usage_1](./resources/stto_usage_1.png)

### Usage 2:
![stto_usage_2](./resources/stto_usage_2.png)

## 🚀 About Me
I'm a Tech enthusiast and a hobby programmer.
You can visit my [Github profile](https://github.com/mainak55512) to see my other works :)

