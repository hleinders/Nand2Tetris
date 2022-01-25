// Adds 1 + ... + 100
@i
M=1
 // i=1
@sum
// another comment and some empty lines


M=0
 // sum=0
(LOOP)
@i
D=M
 // D=i
@100
// comp with whitespaces
D = D - A // D=i-100
@END
D;JGT // if (i-100)>0 goto END
// unwanted indent
  @i
D= M
 // D=i
@sum
M=D+M // sum=sum+i
@i
M=M+1 // i=i+1
@LOOP
0;JMP // goto LOOP
(END)
@END
0;JMP // infinite loop
