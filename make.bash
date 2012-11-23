mkdir -p bin

# build the barbu runner
go build -o bin/runner ./runner/

go build -o bin/server ./net/server
go build -o bin/client ./net/client


# build our ais
go build -o bin/jonai ./jonai
g++ -O2 -o bin/simpleai ./cc_ai_base/*.cc ./simpleai/*.cc
g++ -O2 -o bin/davidai ./cc_ai_base/*.cc ./davidai/*.cc
