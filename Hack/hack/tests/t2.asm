// calc the first N fibonacci numbers
// i: counter
// N: fibonacci numbers to calc
// LAST: address of last number
//
//  1. Prepare setup
//
    @i      // reserve counter space
    @20
    D=A
    @N
    M=D     // N = 20
    @LAST   // reserve space - next mem location is first number
    A=A+1   // first number: 0
    M=0     
    A=A+1   // second number: 1
    M=1
    D=A     // remember addr
    @LAST
    M=D     // and save it to last
    @2
    D=A
    @i
    M=D     // i = 2
