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

//
// 2. Main loop
//
(LOOP)
    @LAST   
    A=M     // set A to last addr
    D=M     // first operand
    A=A-1   // second operand
    D=D+M   // f(n+1) = f(n) + f(n-1)
    @LAST   // get LAST
    A=M     // last addr
    A=A+1   // move forward by 1
    M=D     // write next number
    D=A     // remember A
    @LAST
    M=D     // and update LAST
    @i      // increment counter
    DM=M+1  // save counter and remember val in D
    @N
    D=M-D   // test N-I
    @LOOP
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
(LOOP1)
    @LAST   
    A=M     // set A to last addr
    D=M     // first operand
    A=A-1   // second operand
    D=D+M   // f(n+1) = f(n) + f(n-1)
    @LAST   // get LAST
    A=M     // last addr
    A=A+1   // move forward by 1
    M=D     // write next number
    D=A     // remember A
    @LAST
    M=D     // and update LAST
    @i      // increment counter
    DM=M+1  // save counter and remember val in D
    @N
    D=M-D   // test N-I
    @LOOP
    D;JGE   // jump if >= 0
(END)
    BRK
