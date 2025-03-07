# How to run

```bash
sh run.sh filepath_ratings filepath_movies
```

Tested with `awk` and `split` from https://www.gnu.org/software/coreutils/, and `SQLite` version 3.43.2.

__Important note__: awk and split are aliased as `gawk` and `gsplit` on my machine. Change accordingly inside `run.sh`.


# Result

```
File split (gnu split):

real    0m0.398s
user    0m0.002s
sys     0m0.178s

Inner join (sqlite):

real    0m6.993s
user    0m51.856s (sum of the time of each individual process spawned)
sys     0m1.247s

Average rating (gnu awk):

real    0m0.155s
user    0m0.036s
sys     0m0.005s
```

| GENRE              | RATING   |
| ------------------ | -------- |
| (no genres listed) | 3.326379 |
| Action             | 3.466592 |
| Adventure          | 3.517445 |
| Animation          | 3.614946 |
| Children           | 3.432507 |
| Comedy             | 3.423993 |
| Crime              | 3.685044 |
| Documentary        | 3.705281 |
| Drama              | 3.677185 |
| Fantasy            | 3.511589 |
| Film-Noir          | 3.925728 |
| Horror             | 3.293563 |
| IMAX               | 3.603712 |
| Musical            | 3.554716 |
| Mystery            | 3.670169 |
| Romance            | 3.542712 |
| Sci-Fi             | 3.478143 |
| Thriller           | 3.522964 |
| War                | 3.791466 |
| Western            | 3.585755 |

