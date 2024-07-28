# STTO

Command-line utility written in Go to check total line of code in a file present in a directory.



## Authors

- [Mainak Bhattacharjee](https://github.com/mainak55512)


## Dependencies

- go 1.22.5
- github.com/mattn/go-runewidth v0.0.9
- github.com/olekukonko/tablewriter v0.0.5


## Usage
```bash
❯ stto
+--------------------+------------+-----------------+-----+------+
|     FILE TYPE      | FILE COUNT | NUMBER OF LINES | GAP | CODE |
+--------------------+------------+-----------------+-----+------+
| vim_tutorial.c.swp |          1 |               3 |   0 |    3 |
| c                  |         46 |            1113 |  23 | 1090 |
| py                 |         14 |             165 |   6 |  159 |
| class              |          1 |               8 |   0 |    8 |
| m                  |          1 |               4 |   0 |    4 |
| out                |          1 |               6 |   0 |    6 |
| js                 |          1 |              21 |   2 |   19 |
| java               |          1 |              21 |   3 |   18 |
| css                |          1 |              14 |   0 |   14 |
| html               |          1 |              13 |   0 |   13 |
| cbl                |          1 |              10 |   0 |   10 |
| jl                 |          2 |              16 |   1 |   15 |
| txt                |          1 |              19 |   3 |   16 |
+--------------------+------------+-----------------+-----+------+
Total files:  72
Total lines:  1413
Total gaps:  38
Total code:  1375
```
## 🚀 About Me
I'm a Tech enthusiast and a hobby programmer.
You can visit my [Github profile](https://github.com/mainak55512) to see my other works :)

