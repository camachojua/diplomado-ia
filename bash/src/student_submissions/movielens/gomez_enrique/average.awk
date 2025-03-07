BEGIN {
    FS=","
}

{                                 # skip header column
    split($1, genres, "|")        # split genres (2nd column)
    for(idx in genres) {
        key = genres[idx]
        rating[key] += $2         # sum ratings (3rd column)
        observations[key] += $3   # sum observations (4th column)
    }
}

END {
    for(key in rating){
        printf("%s,%f\n", key, rating[key]/observations[key])
    }
}
