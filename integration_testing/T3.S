	.adra begin
	jmp		ra

%: \fn byteToHexCh
%: \a0 Byte value to convert to hexadecimal representation
%: \a1 Target address of the hexadecimal representation
%: \ds Converts the LSB of the argument 0 into a textual
%: \ds hexadecimal representation which is stored in memory
%: \ds at a location as specified per argument 1. Argument 1
%: \ds must point to a buffer of at least size 2. 
.l byteToHexCh
	push	ro	ra
	push	ro	rb
	
	.adrb	_hexTable
	
	%: get hi-nibble into ri 
	%: and lo-nibble into
	byt		ri
	
	ldca	0x04
	mov		rk	ri
	shr		ri	ra
	
	ldca	0x1C
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
	inc		rj
	stb		rj	rm
	
	pop		rb	ro
	pop		ra	ro
	
	ret		ro
	
.l _hexTable
	strdata MDEyMzQ1Njc4OUFCQ0RFRg==
	
.l foobar
	nop
	
.l begin
	ldca	0xf78
	mov		ri	ra
	
	.adrb	foobar
	mov		rj	rb
	
	.adra	byteToHexCh
	call	ra	ro
	
	ldb		re	rb
	inc		rb
	ldb		rf	rb
	
	ldca	0x12345678
	out		ra	re
	out		ra	rf
	
	.adrc	err
	ldca	0x37
	
	jne		ra	re
	inc		ra
	jne		ra	rf
	
	hlt
	
.l	err
	fail