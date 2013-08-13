//# include "godb.h"
//# include "godb_database.h"

#include <assert.h>
#include <stdint.h>

//#include <sys/types.h>
//#include <string.h>

//#include <stdio.h>

//#include <sys/stat.h>
//#include <sys/param.h>

//#include <sys/uio.h>
//#include <sys/mman.h>
//#ifdef HAVE_SYS_FILE_H
//#include <sys/file.h>
//#endif
//#include <fcntl.h>

//#include <assert.h>
#include <errno.h>
//#include <limits.h>
//#include <stddef.h>
//#include <inttypes.h>
//#include <stdio.h>

//#include <string.h>
////#include <time.h>
#include <unistd.h>

#include "godb.h"
#include "godb_database.h"

	/** The database. */
struct GODB_database {
	HANDLE			header_file;		/** < Handle for writing to the header index file */
	HANDLE			data_file;			/** < Handle for writing to the data file */
	uint32_t		flags;				/** < @ref Environment flags */
	unsigned int 	page_size;			/** < Size of a page, from #GET_PAGESIZE */
	pid_t			pid;				/** < Process ID of the database */
	char 			*file_path;			/** < Path to the db file */
	char 			*data_map;			/** < Memory Map of the data file */
	size_t			data_size;			/** < Size of the data memory map */
	char			*header_map;		/** < Memory Map of the header file */
	size_t 			header_size;		/** < Size of the header memory map */
	//GODB_meta		*header_pages[2];	/** < Pointers to the two header pages */
	//pgno_t			max_pages;			/** < map_size / page_size */
	//GODB_dbx		db_info;			/** < Array of static DB info */

	// Generations
	// Indexes
	// Memory Map(s)
	// Buffer(s)
};

HANDLE godb_open_file(char *path, int file_mode) {
	//int error;
	HANDLE file_handle;
	//int received;
	off_t file_size;//, rsize;

	file_handle = open(path, READ_WRITE_CREATE, file_mode);

	// Unable to open file handle
	if (file_handle == INVALID_HANDLE_VALUE) {
		return INVALID_HANDLE_VALUE;
		//error = ErrCode(); // Line 3620 in mdb.c
		//goto fail_with_error;
	}

	// Seek from beginning to end to find file size
	file_size = lseek(file_handle, 0, SEEK_END);
	if (file_size == -1) {
		return INVALID_HANDLE_VALUE;
		//goto fail_with_error;
	}

	return GODB_SUCCESS;

	//fail_with_error:
	//	return INVALID_HANDLE_VALUE;
}

/*
int
godb_get_flags(GODB_database *db, godb_flags *ret_flags) {
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


/** Read the Header data from the Database
 * @param[in]	db: 				Database handle
 * @param[out]	header_data:		Byte[] buffer to hold raw header data
 * @return 							0 on success, non-zero on failure
 */
 // -> mdb_env_read_header
 /*
static int
godb_read_header(GODB_database *env, void *header_data) {
	//GODB_pagebuf pbuf;
	int header_filehandle, received, offset;
	int page_size = 1024;

	received = pread(header_filehandle, header_data, page_size, offset);

	if (received != page_size) { // Reached the end of the file
		if (received == 0 && offset == 0) { // No data in the file
			return ENOENT;
		}
		// If the response is less than zero it means there was an error
		//received = received < 0 ? (int) ErrCode() : GODB_INVALID;
		//DPRINTF("Read header with result: %s", godb_error_string(received));
		return received;
	}

	return 0;
}
*/

// Initialize new DB environment
// -> mdb_env_init_meta

/** Initialize an enviornment variable
 * @param[out]	out_env:			Created environment handle
 * @return 							0 on success, non-zero on failure
 */
// -> mdb_env_create
 /*
int godb_create(GODB_database **out_db) {
	GODB_database  *db;

	// Allocate memory for database variable
	db = calloc(1, sizeof(GODB_database));
	// Report an error if calloc failed
	if (!db) return ENOMEM;

	// Remove next two lines
	char *buffer;
	godb_read_header(db, &buffer);
	//create_file(buffer, 0666, 0)
	// Hydrate database

	// Assign the new database into the return
	*out_db = db;
	return GODB_SUCCESS;
}
*/


/** mdb_put */
/*
int
godb_put(GODB_database db, uint64_t key, void *data) {
//mdb_put(MDB_txn *txn, MDB_dbi dbi, MDB_val *key, MDB_val *data, unsigned int flags)
	assert(data != NULL);

	// MDB_cursor mc;
	//MDB_xcursor mx;

	// Check to make sure the provided flags don't prevent data writing
	//if ((flags & (MDB_NOOVERWRITE|MDB_NODUPDATA|MDB_RESERVE|MDB_APPEND|MDB_APPENDDUP)) != flags)
	//	return EINVAL;

	//mdb_cursor_init(&mc, txn, dbi, &mx);
	//return mdb_cursor_put(&mc, key, data, flags);

	return -1;
}
*/




//void write_to_file(char *filePath, void* data) {

//}