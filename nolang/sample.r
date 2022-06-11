repeat y 10 @y
ret

:y
    repeat x 10 @x
    ret

:x
    print
        concat
            concat "X: " str x
            concat " Y: " str y
    ret