set a arr-new

set a arr-add a "a"
set a arr-add a "b"
set a arr-add a "c"
set a arr-add a "d"
set a arr-add a "e"

set aptr 0
set val 0

set x 100
:lp
    set x sub x 1
    call @get-val
    print val
    sleep 0.01
    loop x @lp

:get-val
    set aptr add aptr 1
    set aptr mod aptr arr-len a
    set val arr-get a aptr
    ret