.l std::test::memcmp_buf1
	strdata YWJjZGU=
.l std::test::memcmp_buf2
	strdata YWJjZGU=
.l std::test::memcmp_buf3
	strdata MTIzNDU=
	
.x std::test::memcmp
	push	ro	rc
	
	.adrc	std::test::memcmp_buf1
	mov		re	rc
	
	.adrc	std::test::memcmp_buf2
	mov		rf	rc
	
	.adrc	std::test::memcmp_buf3
	mov		rg	rc
	
	.xadrc	std::memcmp
	mov		rh	rc
	
	.adrc	std::test::memcmp_fail
	
	mov		ri	re
	mov		rj	rf
	call	rh	ro
	jiz		rc	rn
	
	mov		ri	re
	mov		rj	rg
	call	rh	ro
	jnz		rc	rn
	
	pop		rc	ro
	ret		ro
	
.l std::test::memcmp_fail
	fail

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