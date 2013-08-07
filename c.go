package eventstore

/*
//# include <assert.h>
//# include <stdint.h>
//# include <string.h>
//# include <stdio.h>
//# include <stdlib.h>

void Copy_memory_loop(void* dest, void* src, int length) {
	unsigned char* cdest = (unsigned char*) dest;
	unsigned char* csrc = (unsigned char*) src;

	int i;
	for (i = 0; i < length; i++) {
		cdest[i] = csrc[i];
	}
}

void write_memory_rep_stosq(void* dest, void* src, int length) {
	//unsigned char* cdest = (unsigned char*) dest;
	//unsigned char* csrc = (unsigned char*) src;

	//memset(cdest, csrc, 8);
  asm("cld\n"
      "rep stosq"
      : : "D" (dest), "c" (1), "a" (src) );
}

*/
import "C"
