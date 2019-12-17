.defm xcallc .xadrc #0 \ .call rc

.l init
	.adrc	stack
	mov		ro	rc
	
.l std::test::run_tests
.m	xcallc	std::test::memcpy
	
	.xadrc	std::test::memcmp
	.call	rc
	
.l std::test::run_tests_end
	hlt
	
.l stack
	.alloc	256