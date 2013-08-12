#include "godb.h"

int temp(int value) {
	return value * 10;
}

char *godb_version(int *major, int *minor, int *patch)
{
	if (major) *major = GODB_VERSION_MAJOR;
	if (minor) *minor = GODB_VERSION_MINOR;
	if (patch) *patch = GODB_VERSION_PATCH;
	return GODB_VERSION_STRING;
}
/*


#ifndef _GNU_SOURCE
#define _GNU_SOURCE 1
#endif
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/param.h>

#include <sys/uio.h>
#include <sys/mman.h>
#ifdef HAVE_SYS_FILE_H
#include <sys/file.h>
#endif
#include <fcntl.h>

#include <assert.h>
#include <errno.h>
#include <limits.h>
#include <stddef.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
//#include <time.h>
#include <unistd.h>


#if defined(__APPLE__) || defined (BSD) 
# define GODB_FDATASYNC		fsync
#elif defined(ANDROID)
# define GODB_FDATASYNC		fsync
#endif

#ifdef USE_VALGRIND
#include <valgrind/memcheck.h>
#define VGMEMP_CREATE(h,r,z)    VALGRIND_CREATE_MEMPOOL(h,r,z)
#define VGMEMP_ALLOC(h,a,s) VALGRIND_MEMPOOL_ALLOC(h,a,s)
#define VGMEMP_FREE(h,a) VALGRIND_MEMPOOL_FREE(h,a)
#define VGMEMP_DESTROY(h)	VALGRIND_DESTROY_MEMPOOL(h)
#define VGMEMP_DEFINED(a,s)	VALGRIND_MAKE_MEM_DEFINED(a,s)
#else
#define VGMEMP_CREATE(h,r,z)
#define VGMEMP_ALLOC(h,a,s)
#define VGMEMP_FREE(h,a)
#define VGMEMP_DESTROY(h)
#define VGMEMP_DEFINED(a,s)
#endif


*/


/** Copy input string to the output string */
/** http://stackoverflow.com/questions/8359966/strdup-returning-address-out-of-bounds */
/*
char *dupstr(const char *s) {
    char *const result = malloc(strlen(s) + 1);
    if (result != NULL) {
        strcpy(result, s);
    }
    return result;
}
*/
/*

int temp(int temp) {
	int i = temp * 3;
	return i;
}


*/


//#define BLOCK_SIZE	32 // 32 bytes per block
	/* Max bytes to write in one call */
//#define MAX_WRITE		(0x80000000U >> (sizeof(ssize_t) == 4))

//#define GODB_ERRCODE_ROFS	EROFS
//#define GODB_CLOEXEC		O_CLOEXEC

//#define READ_WRITE_CREATE	O_RDWR | O_CREAT | GODB_CLOEXEC


/** Return the library version info. */
/*char *godb_version(int *major, int *minor, int *patch)
{
	if (major) *major = GODB_VERSION_MAJOR;
	if (minor) *minor = GODB_VERSION_MINOR;
	if (patch) *patch = GODB_VERSION_PATCH;
	return GODB_VERSION_STRING;
}*/
	/** The name of the data file in the DB enviornment */
//#define HEADER_FILE_NAME "/header.godb"
	/** The name of the header file in the DB environment */
//#define DATA_FILE_NAME	"/data.godb"
	/** The number of bytes to offset to account for File Info when parsing the files */
//#define FILE_INFO_SIZE	8


//#include "godb_database.h"
//#include "godb_database.c"



	 /** Info written to the top of each GODB file */
//typedef struct GODB_file_info {
//	int 			bin_version;		/** < Version number of the compiled binary that last managed the file */
//	int 			format_version;		/** < Version number of the file format to detect legacy files */
//} GODB_file_header;

// How to keep these alligned on disk?
//typedef struct GODB_aggregate_header {
//	uint64_t 	id;					/** < Unique 64 bit identifer for an aggregate */
//} GODB_aggregate_header;

//typedef struct GODB_event_header {
//	uint32_t 	length;				/** < Length (in bytes) of the event data */
//	uint32_t 	event_type;			/** < Client supplied type identifier to allow deserializing of events */
//} GODB_event_header;


/*
int
godb_set_flags(GODB_database *db, unsigned int flags, unsigned short turnOn) {
	// if ((flag & CHANGEABLE) != flag) return EINVAL;

	if (turnOn) {
		db->flags |= flags;
	} else {
		db->flags &= ~flags;
	}
	return GODB_SUCCESS;
}

int
godb_get_flags(GODB_database *db, unsigned int *ret_flags) {
	if (!db || !ret_flags) { return EINVAL; }

	*ret_flags = db->flags;
	return GODB_SUCCESS;
}

int
godb_get_path(GODB_database *db, const char **ret_path) {
	if (!db || !ret_path) { return EINVAL; }

	*ret_path = db->file_path;
	return GODB_SUCCESS;
}
*/







