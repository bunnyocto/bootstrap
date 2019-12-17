.l init
	.adrc	stack
	mov		ro	rc
	
.l std::test::run_tests
	.xadrc	std::test::memcpy
	.call	rc
	
	.xadrc	std::test::memcmp
	.call	rc
	
.l std::test::run_tests_end
	hlt
	
.l stack
	.alloc	256