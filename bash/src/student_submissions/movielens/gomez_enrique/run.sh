filepath_ratings=$1
filepath_movies=$2

echo "File split:"
time gsplit -d -C 50M $filepath_ratings --additional-suffix .csv "tmp_split_"
tail -n +2 tmp_split_00.csv > tmp_split_00.tmp && mv tmp_split_00.tmp tmp_split_00.csv

for filepath_split in tmp_split_*.csv; do
    filename=$(basename -- "$filepath_split")
    extension="${filename##*.}"
    filename="${filename%.*}"
    sqlite3 ":memory:" <<SQL &
.output /dev/null
pragma journal_mode = OFF;
pragma synchronous = OFF;
pragma journal_size_limit = 0;

.import $filepath_movies movies --csv
create unique index idx_movieid on movies(movieId);

create table ratings(userId, movieId, rating, timestamp);
.import $filepath_split ratings --csv

.mode csv
.header off
.out tmp_innerjoin_$filename.csv

select genres, sum(rating), count(rating)
from movies inner join ratings
on movies.movieId = ratings.movieId
group by genres;

.exit
SQL
done

echo "\nInner join:"
time wait

echo "\nAverage rating:"
cat tmp_innerjoin_*.csv > tmp_ratings_min.csv
time gawk -f average.awk tmp_ratings_min.csv | sort > tmp_out.csv

cat tmp_out.csv

rm tmp_*.csv
