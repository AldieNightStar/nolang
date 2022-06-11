set a arr-all arr-new
    10 20 30 40 50 60 70
!!

set cnt arr-len a
:foreach
    set id sub arr-len a cnt
    set val arr-get a id
    print val
    loop cnt @foreach