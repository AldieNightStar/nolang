repeat y 10000 @y
ret

:y
    repeat x 10000 @x
    ret

:x
    if == mod x 100 0 @print-xy
    sleep 0.00001
    ret

:print-xy
    print
        concat
            concat "X: " str x
            concat " Y: " str y
    ret