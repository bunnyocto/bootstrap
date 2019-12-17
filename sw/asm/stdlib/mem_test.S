.l std::test::memcpy_origbuf
    strdata MTIzNDU=
.l std::test::memcpy_cpybuf
	strdata MDAwMDA=
	
.x std::test::memcpy
	push	ro	rc
	push	ro	rb
	push	ro	ra

	.adra	std::test::memcpy_cpybuf
	.adrb	std::test::memcpy_origbuf
	
	mov		ri	ra
	mov		rj	rb
	ldp		rk	5
	
	.xadrc	std::memcpy
	call	rc	ro
	
	mov		ri	ra
	mov		rj	rb
	ldp		rk	5

	.xadrc	std::memcmp
	call	rc	ro
	.adrc	std::test::memcpy_mismatch
	
	jiz		rc	rn
	
	pop		ro	ra
	pop		ro	rb
	pop		ro	rc
	ret		ro
	
.l std::test::memcpy_mismatch
	fail