	%: Jump over data declerations
	.adra begin
	jmp ra
	
	%: This is the "Hello, world!" string (b64)
.l hw
	strdata SGVsbG8sIHdvcmxkIQ==

.l begin
	%: 0x12345678 is the debug output console port
	%: we'll load it into re
	ldca 0x12345678
	mov re ra

	%: load the address of hw into rd
	.adra hw
	mov rd ra

	%: length of the string (hw)
	ldcb 0x0D

	%: load the loop address into ra
	.adra loop
.l loop

	%: fetch byte, inc, dec, repeat
	ldb rc rd
	out re rc
	inc rd
	dec rb
	jnz ra rb

	hlt %: death