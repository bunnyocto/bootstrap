%: NAME std::memcpy
%: ARGS
%:  ri  Destination
%:  rj  Source
%:  rk  N
%: DESC Copies N bytes from Source into Destination. The
%:      regions must not overlap.
%: RET  Nothing
.x std::memcpy
	.var	Dst	ri
	.var	Src	rj
	.var	N	rk
	.var	Tmp	rl

	.save	rc
	
	.whilenz N
		ldb		Tmp	Src
		stb		Dst	Tmp
		
		inc		Dst
		inc		Src
		dec		N
	.end

	.ret
	
.eof

%: NAME std::memcmp
%: ARGS
%:  ri A
%:  ri B
%:  rk N
%: DESC Compares memory regions. At most N bytes will be compared.
%:      The regions must not overlap. 
%: RET 0 if regions are not equal, otherwise `memcmp` returns 1.
.x std::memcmp
	.var	A	ri
	.var	B	rj
	.var	N	rk
	.var	Ta	rl
	.var	Tb	rm

	.save	rc
	
	.whilenz N
		ldb		Ta	A
		ldb		Tb	B
		
		.if 	ne	Ta	Tb
			ldp		rn	0
			.ret
		.end
		
		inc		A
		inc		B
		dec		N
	.end

	ldp		rn	1
	.ret
	
.eof