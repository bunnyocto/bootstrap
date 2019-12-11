%: \fn byteToHexCh
%: \a0 Byte value to convert to hexadecimal representation
%: \a1 Target address of the hexadecimal representation
%: \ds Converts the LSB of the argument 0 into a textual
%: \ds hexadecimal representation which is stored in memory
%: \ds at a location as specified per argument 1. Argument 1
%: \ds must point to a buffer of at least size 2. 
.l _hexTable
	strdata MDEyMzQ1Njc4OUFCQ0RFRg==
.l byteToHexCh
	push	sp	ra
	push	sp	rb
	
	ldca	0x04
	.adrb	_hexTable
	
	%: get hi-nibble into ri 
	%: and lo-nibble into rk
	mov		rk	ri
	shr		ri	ra
	
	shl		rk	ra
	shr		rk	ra
	
	%: fetch entries from _hexTable and write to 
	%: target address. 
	add		rb	ri
	ldb		rm	rb
	sub		rb	ri
	stb		rj	rm
	
	add		rb	rk
	ldb		rm	rb
	stb		rj	rb
	
	pop		ra	sp
	pop		rb	sp
	
	ret