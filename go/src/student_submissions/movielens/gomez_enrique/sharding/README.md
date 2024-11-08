# How to run

```console
$ go build
$ ./movielens file_to_split total_splits out_folder file_to_merge
```

where:
- `file_to_split`: ratings.csv file
- `total_splits`: total number of files to be splitted into
- `out_folder`: directory where the splitted files will be stored
- `file_to_merge`: movies.csv file

# Example

```bash
./movielens "ratings.csv" 10  "." "movies.csv"
```

## Results

__Split: 229ms__

__Merge & Count: 10363ms__

| GENRE       | AVERAGE_RATING                   |
|-------------|----------------------------------|
| (no genres listed)         | 3.326379239118188 |
| Action      | 3.466591472228235                |
| Adventure   | 3.517444379462875                |
| Animation   | 3.614946348438093                |
| Children    | 3.4325074920278045               |
| Comedy      | 3.423992522260525                |
| Crime       | 3.6850431095379377               |
| Documentary | 3.7052805249822454               |
| Drama       | 3.6771844525139366               |
| Fantasy     | 3.5115889157486                  |
| Film-Noir   | 3.9257258540768367               |
| Horror      | 3.2935633075659174               |
| IMAX        | 3.6037121959523324               |
| Musical     | 3.5547170809260242               |
| Mystery     | 3.670169244577933                |
| Romance     | 3.542711629764427                |
| Sci-Fi      | 3.4781437345156516               |
| Thriller    | 3.522964338694794                |
| War         | 3.7914657875591984               |
| Western     | 3.585752382527443                |

