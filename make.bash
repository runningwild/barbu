mkdir -p bin

# build the barbu runner
go build -o barbu .

# build our ais
go build -o bin/jonai ./jonai
g++ -o bin/davidai ./davidai/*.cc
