%: NAME std::memcpy
%: ARGS
%:  ri  Destination
%:  rj  Source
%:  rk  N
%: DESC Copies N bytes from Source into Destination. The
%:      regions must not overlap.
%: RET  Nothing
.x std::memcpy
	push	ro	rc
	push	ro	rb
	
	.adrc	std::memcpy_loop_end
	.adrb	std::memcpy_loop
.l std::memcpy_loop
	jiz		rc	rk

	ldb		rl	rj
	stb		ri	rl
	
	inc		ri
	inc		rj
	dec		rk
	jmp		rb

.l std::memcpy_loop_end
	pop		rb	ro
	pop		rc	ro
	ret		ro

%: NAME std::memcmp
%: ARGS
%:  ri A
%:  ri B
%:  rk N
%: DESC Compares memory regions. At most N bytes will be compared.
%:      The regions must not overlap. 
%: RET 0 if regions are not equal, otherwise `memcmp` returns 1.
.x std::memcmp
	push	ro	rc
	push	ro	rd
	
	.adrc	std::memcmp_loop_end
	mov		rd	rc
	.adrc	std::memcmp_loop_end_false
	.adrb	std::memcmp_loop
	
.l std::memcmp_loop
	jiz		rd	rk

	ldb		rl	rj
	ldb		rm	ri
	
	jne		rl	rm
	
	inc		ri
	inc		rj
	dec		rk
	jmp		rb

.l std::memcmp_loop_end
	pop		rd	ro
	pop		rc	ro
	ldp		rn	1
	ret		ro
	
.l std::memcmp_loop_end_false
	pop		rd	ro
	pop		rc	ro
	ldp		rn	0
	ret		ro