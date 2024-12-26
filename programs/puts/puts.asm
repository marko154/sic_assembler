puts	START   0
	CLEAR   X
	LDA	expr
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
expr	EQU     3 *  (len - 1)
nlA	WORD	0
	END     puts