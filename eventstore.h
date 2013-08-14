#ifndef _ES_H_
#define _ES_H_

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
	
#define ES_VERSION_MAJOR	0						/** Library major version */
#define ES_VERSION_MINOR	1 						/** Library minor version */
#define ES_VERSION_PATCH	0 						/** Library patch version */
#define ES_VERSION_DATE	"January 10, 2013" 			/** The release date of this library version */

	/** Combine args a,b,c into a single integer for easy version comparisons */
#define ES_VERSION_INT(a,b,c)	(((a) << 24) | ((b) << 16) | (c))
	/** A stringifier for the version info */
#define ES_VERSION_STR(a,b,c,d)	"ES " #a "." #b "." #c ": (" d ")"

	/** The full library version as a single integer */
#define ES_VERSION_FULL														\
	ES_VERSION_INT(															\
		ES_VERSION_MAJOR,													\
		ES_VERSION_MINOR,													\
		ES_VERSION_PATCH);											

	/** The full library version as a C string */
#define	ES_VERSION_STRING													\
	ES_VERSION_STR(															\
		ES_VERSION_MAJOR,													\
		ES_VERSION_MINOR,													\
		ES_VERSION_PATCH,													\
		ES_VERSION_DATE);

char *es_version(int *major, int *minor, int *patch);		/** Return the library version info. */

/** @ defgroup errors 	Return Codes
*
*	Avoid conflict with BerkelyDB (-30800 to -30999) and MDB (-30799 to -30783)
*	@{
*/
#define ES_SUCCESS 				0 							/** Successful result */
#define ES_ERROR				(-30600)					/** Generic error */
#define ES_NOTFOUND 			(ES_ERROR - 1) 		/** Key not found during get (EOF) */
#define ES_CORRUPTED 			(ES_ERROR - 2)		/** CheckSum failed */
#define ES_PANIC 				(ES_ERROR - 3) 		/** Update of meta page failed, probably I/O error */
#define ES_VERSION_MISMATCH 	(ES_ERROR - 4) 		/** Environment version mismatch */
#define ES_INVALID 				(ES_ERROR - 5)		/** Invalid ES file */
#define ES_MAP_FULL				(ES_ERROR - 6)		/** Environment mapsize reached */
#define ES_PAGE_FULL			(ES_ERROR - 7)		/** Page ran out of space - internal error */
#define ES_MAP_RESIZED			(ES_ERROR - 8)		/** Database contents grew benyond environment mapsize */
#define ES_INCOMPATIBLE			(ES_ERROR - 9)		/** Database flags changes (or would change) */
/** @} */


void es_open(char* path);






#endif /* _ES_H_ */

