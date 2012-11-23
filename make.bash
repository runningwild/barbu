mkdir -p bin

# build the barbu runner
go build -o bin/runner ./runner/

# build the net player binary
go build -o ./bin/net_player ./net_player/

# build our ais
go build -o bin/jonai ./jonai
g++ -O2 -o bin/davidai ./davidai/*.cc
