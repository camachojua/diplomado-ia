#Correr el comando: ./run.sh en terminal de MAC. 
#Colocar el archivo de movies y rattings en la carpeta: /./ADCC_project/csv_files/
export GO111MODULE=on
echo "GO111MODULE=on"
go build
echo "go build OK"
go run movielens.go