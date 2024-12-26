puts	START   0
	CLEAR   X
loop	LDCH    txt, X
	JSUB    putc
	TIX     #len
	JLT     loop
halt	J       halt


.print char in register A
putc	WD      #1
	RSUB

nl	STA     nlA
	LDCH    #10
	JSUB    putc
	LDA     nlA
	RSUB

txt	BYTE    C'hello world'
end	EQU     *
len	EQU     end - txt
nlA	WORD	0
	END     puts