.l std::test::run_tests
	.xadrc	std::test::memcpy
	call	rc	ro
	
	.xadrc	std::test::memcmp
	call	rc	ro
	
.l std::test::run_tests_end
	hlt