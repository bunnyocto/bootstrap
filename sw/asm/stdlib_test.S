.defm xcallc .xadrc #0 \ .call rc

.l init
	.adrc	stack
	mov		ro	rc
	
.l std::test::run_tests
.m	xcallc	std::test::memcpy
.m	xcallc	std::test::memcmp
.m	xcallc	std::test::memset
	
.l std::test::run_tests_end
	hlt
	
.l stack
	.alloc	256