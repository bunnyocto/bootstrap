%: NAME std::memcpy
%: ARGS
%:  ri  Destination
%:  rj  Source
%:  rk  N
%: RET  Nothing
.x std::memcpy
	push	ro	rc
	
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
	pop		rc	ro
	ret		ro