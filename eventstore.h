#ifndef _EVENTSTORE_H_
#define _EVENTSTORE_H_

#include <stdarg.h> /* Needed for the definition of va_list */
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include <unistd.h>

#include <sys/uio.h>
#ifdef HAVE_SYS_FILE_H
#include <sys/file.h>
#endif
#include <fcntl.h>

#include <sys/mman.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <errno.h>

/*

#include <sys/types.h>
#include <sys/stat.h>
#include <sys/param.h>

#include <assert.h>
#include <errno.h>
#include <limits.h>
#include <stddef.h>
#include <inttypes.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

*/
	
#define EVENTSTORE_VERSION_MAJOR	0						/** Library major version */
#define EVENTSTORE_VERSION_MINOR	1 						/** Library minor version */
#define EVENTSTORE_VERSION_PATCH	0 						/** Library patch version */
#define EVENTSTORE_VERSION_DATE	"January 10, 2013" 			/** The release date of this library version */

	/** Combine args a,b,c into a single integer for easy version comparisons */
#define EVENTSTORE_VERSION_INT(a,b,c)	(((a) << 24) | ((b) << 16) | (c))
	/** A stringifier for the version info */
#define EVENTSTORE_VERSION_STR(a,b,c,d)	"EVENTSTORE " #a "." #b "." #c ": (" d ")"

	/** The full library version as a single integer */
#define EVENTSTORE_VERSION_FULL														\
	EVENTSTORE_VERSION_INT(															\
		EVENTSTORE_VERSION_MAJOR,													\
		EVENTSTORE_VERSION_MINOR,													\
		EVENTSTORE_VERSION_PATCH);											

	/** The full library version as a C string */
#define	EVENTSTORE_VERSION_STRING													\
	EVENTSTORE_VERSION_STR(															\
		EVENTSTORE_VERSION_MAJOR,													\
		EVENTSTORE_VERSION_MINOR,													\
		EVENTSTORE_VERSION_PATCH,													\
		EVENTSTORE_VERSION_DATE);

char *eventstore_version(int *major, int *minor, int *patch);		/** Return the library version info. */

/** @ defgroup errors 	Return Codes
*
*	Avoid conflict with BerkelyDB (-30800 to -30999) and MDB (-30799 to -30783)
*	@{
*/
#define EVENTSTORE_SUCCESS 				0 							/** Successful result */
#define EVENTSTORE_ERROR				(-30600)					/** Generic error */
#define EVENTSTORE_NOTFOUND 			(EVENTSTORE_ERROR - 1) 		/** Key not found during get (EOF) */
#define EVENTSTORE_CORRUPTED 			(EVENTSTORE_ERROR - 2)		/** CheckSum failed */
#define EVENTSTORE_PANIC 				(EVENTSTORE_ERROR - 3) 		/** Update of meta page failed, probably I/O error */
#define EVENTSTORE_VERSION_MISMATCH 	(EVENTSTORE_ERROR - 4) 		/** Environment version mismatch */
#define EVENTSTORE_INVALID 				(EVENTSTORE_ERROR - 5)		/** Invalid EVENTSTORE file */
#define EVENTSTORE_MAP_FULL				(EVENTSTORE_ERROR - 6)		/** Environment mapsize reached */
#define EVENTSTORE_PAGE_FULL			(EVENTSTORE_ERROR - 7)		/** Page ran out of space - internal error */
#define EVENTSTORE_MAP_RESIZED			(EVENTSTORE_ERROR - 8)		/** Database contents grew benyond environment mapsize */
#define EVENTSTORE_INCOMPATIBLE			(EVENTSTORE_ERROR - 9)		/** Database flags changes (or would change) */
/** @} */


void eventstore_open(char* path);






#endif /* _EVENTSTORE_H_ */

