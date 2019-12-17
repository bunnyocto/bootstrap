.l std::test::memcmp_buf1
	strdata YWJjZGU=
.l std::test::memcmp_buf2
	strdata YWJjZGU=
.l std::test::memcmp_buf3
	strdata MTIzNDU=
	
.x std::test::memcmp
	.save	rc
	
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
	ldp		rk	5
	call	rh	ro
	jiz		rc	rn
	
	mov		ri	re
	mov		rj	rg
	ldp		rk	5
	call	rh	ro
	jnz		rc	rn
	
	.ret
	
.eof
	
.l std::test::memcmp_fail
	fail

.l std::test::memcpy_origbuf
    strdata MTIzNDU=
.l std::test::memcpy_cpybuf
	strdata MDAwMDA=
	
.x std::test::memcpy
	.save	rc
	.save	rb
	.save	ra

	.adra	std::test::memcpy_cpybuf
	.adrb	std::test::memcpy_origbuf
	
	mov		ri	ra
	mov		rj	rb
	ldp		rk	5
	
	.xadrc	std::memcpy
	.call	rc
	
	mov		ri	ra
	mov		rj	rb
	ldp		rk	5

	.xadrc	std::memcmp
	.call	rc
	.adrc	std::test::memcpy_mismatch
	
	jiz		rc	rn
	
	.ret
	
.l std::test::memcpy_mismatch
	fail
	
.eof