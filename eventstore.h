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

#include <sys/socket.h>

/*

#include <sys/param.h>

#include <assert.h>
#include <limits.h>
#include <stddef.h>
#include <inttypes.h>
#include <time.h>

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
#define ES_SUCCESS 						0 					/** Successful result */
#define ES_ERROR						(-30600)			/** Generic error */

//#define ES_NOTFOUND 					(ES_ERROR - 1) 		/** Key not found during get (EOF) */
//#define ES_CORRUPTED 					(ES_ERROR - 2)		/** CheckSum failed */
//#define ES_PANIC 						(ES_ERROR - 3) 		/** Update of meta page failed, probably I/O error */
//#define ES_VERSION_MISMATCH 			(ES_ERROR - 4) 		/** Environment version mismatch */
//#define ES_MAP_FULL						(ES_ERROR - 5)		/** Environment mapsize reached */
//#define ES_PAGE_FULL					(ES_ERROR - 6)		/** Page ran out of space - internal error */
//#define ES_MAP_RESIZED					(ES_ERROR - 7)		/** Database contents grew benyond environment mapsize */
//#define ES_INCOMPATIBLE					(ES_ERROR - 8)		/** Database flags changes (or would change) */

#define ES_FILE_NOTFOUND 				(ES_ERROR - 100)	/** ES file was not found */
#define ES_FILE_INVALID 				(ES_ERROR - 200)	/** ES file is invalid */

#define ES_SETTINGS_FILE_NOTFOUND 		(ES_FILE_NOTFOUND - 1)	/** ES Settings file was not found */
#define ES_SETTINGS_FILE_INVALID 		(ES_FILE_INVALID - 1)	/** ES Settings file is invalid */
#define ES_HEADER_FILE_NOTFOUND 		(ES_FILE_NOTFOUND - 2)	/** ES Header file was not found */
#define ES_HEADER_FILE_INVALID 			(ES_FILE_INVALID - 2)	/** ES Header file is invalid */
#define ES_GEN_FILE_NOTFOUND 			(ES_FILE_NOTFOUND - 3)	/** ES Generation file was not found */
#define ES_GEN_FILE_INVALID 			(ES_FILE_INVALID - 3)	/** ES Generation file is invalid */
#define ES_DATA_FILE_NOTFOUND 			(ES_FILE_NOTFOUND - 10)	/** ES Data file was not found */
#define ES_DATA_FILE_INVALID 			(ES_FILE_INVALID - 10)	/** ES Data file is invalid */


/** @} */

	/** Represents an opaque handle to the real database struct */
//typedef struct ES_handle ES_handle;

	/** A domain creates a top level partition separating internal headers */
//typedef struct ES_domain ES_domain;

	/** A kind represents a type of Aggregate in the event store */
//typedef struct ES_kind ES_kind;

//ES_handle* es_open(char* path);
//void es_close(ES_handle* writer);
//int es_open(ES_database** database, char* path);

	/** Targeting min workable size that fits cleanly in 4kb */
#define ES_GEN_PAGE_SIZE				4096
#define ES_COMMAND_SIZE					512
#define ES_COMMAND_HEADER_SIZE			32		// 32 is command overhead size
#define ES_MAX_DATA_SIZE				ES_COMMAND_SIZE - ES_COMMAND_HEADER_SIZE
#define ES_COMMANDS_PER_PAGE			ES_GEN_PAGE_SIZE / ES_COMMAND_SIZE

typedef struct ES_writer ES_writer;

typedef struct ES_batch_entry ES_batch_entry;
typedef struct ES_batch ES_batch;

typedef struct ES_put_command ES_put_command;

ES_writer* es_open_write(char* path);
void es_close_write(ES_writer* writer);
ES_batch* es_alloc_batch(ES_writer* writer, 
	uint32_t domain_id, 
	uint32_t kind_id, 
	uint64_t aggregate_id, 
	char count);





#endif /* _ES_H_ */

