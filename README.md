# NoLang

### Just for fun!

# Sample
```r
# Init the logger
set log logger-new "Default Logger"

# Connection
log "Initializing..."
set connection socket-connect "192.168.0.1" 8080
!if-else socket-is-ok connection @begin @failure
ret

# When connection is ok
:begin
    log "Connection retreived"
    log "Sending data..."
    :retry
    socket-send connection to-json user-get-data
    sleep 1 # Sleep 1 second
    set packets socket-receive connection
    !if < arr-len packets 1 @retry
    arr-each packets packet @process-packet
    log "Operation done. Disconnecting..."
    socket-disconnect connection
    log "Done"
ret

# When connection is bad
:failure
    log "Connection fail"
ret

:process-packet
    log concat "Got packet with size: " str arr-len packet
    central-operate packet
    log "Packet is processed!"
ret
```

# API
* Run a program
```go
// Load Scope
scope := nolang.LoadFile("program.r")

// Run Scope
scope.Run()
```
* Add Custom variables
```go
// This allows NoLang to have getter/setter which can control variables outside of it environment

// This function will set count-get as function which returns cnt variable
scope.Mem["count-get"] = nolang.ValueGetter(func() (any, error) { return cnt, nil })

// This function will set count-set as function which will set new value to cnt
scope.Mem["count-set"] = nolang.ValueSetter(func(a any) error { cnt = a.(float64); return nil })
```
* Add custom functions
```go
// To make NoLang call your functions, set your functions with nolang.NoFunc1 .. nolang.NoFunc4
// nolang.NoFunc1 has generic, and default values as first arguments. Last argument is a function

// print twice 10
scope.Mem["twice"] = nolang.NoFunc1[float64](0, func(a float64) (any, error) {
    return a * 2, nil
})

// print name-age "HaxiDenti" 30
scope.Mem["name-age"] = nolang.NoFunc2[string, float64]("", 0, func(name string, age float64) (any, error) {
    return fmt.Sprintf("name: %s and age: %.2f", name, age), nil
})

// print add3 10 20 30
scope.Mem["add3"] = nolang.NoFunc3[float64, float64, float64](0, 0, 0, func(a, b, c float64) (any, error) {
    return a + b + c, nil
})

// print add4 10 20 30 40
scope.Mem["add4"] = nolang.NoFunc4[float64, float64, float64, float64](0, 0, 0, 0, func(a, b, c, d float64) (any, error) {
    return a + b + c + d, nil
})
```
* Advanced functions
```go
// This function is setting execution cursor to the beginning
// Such function can control core of the language, even change memory or call-stack
scope.Mem["reset"] = nolang.NoFunc(func(s *nolang.Scope) (any, error) {
    s.Pos = 0
    return nil, nil
})
```
* Just set up the variables
```go
scope.Mem["name"] = "Ihor"
```
* Get some variables
```go
fmt.Println(scope.Mem["result"])
```