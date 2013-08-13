#include "eventstore.h"

char *eventstore_version(int *major, int *minor, int *patch)
{
	if (major) *major = EVENTSTORE_VERSION_MAJOR;
	if (minor) *minor = EVENTSTORE_VERSION_MINOR;
	if (patch) *patch = EVENTSTORE_VERSION_PATCH;
	return EVENTSTORE_VERSION_STRING;
}

#define SETTINGS_FILE_NAME		"es_settings_file.esdb"
#define HEADER_FILE_NAME 		"es_header_file.esdb"
#define HEADER_SWAP_FILE_NAME	"es_header_swap_file.esdb"
#define DATA_FILE_NAME 			"es_data_file.esdb"
#define DATA_GEN_FILE_NAME		"es_data_gen_file_XXXX.esdb"

#define SETTINGS_FILE_NAME_SIZE			sizeof(SETTINGS_FILE_NAME) - 1
#define HEADER_FILE_NAME_SIZE 			sizeof(HEADER_FILE_NAME) - 1
#define HEADER_SWAP_FILE_NAME_SIZE		sizeof(HEADER_SWAP_FILE_NAME) - 1
#define DATA_FILE_NAME_SIZE				sizeof(DATA_FILE_NAME) - 1
#define DATA_GEN_FILE_NAME_SIZE			sizeof(DATA_GEN_FILE_NAME) - 1

#define EVENTSTORE_FILE_KEY			0xC0DEC0DE		/** Stamp to identify a file as EVENTSTORE valid and check for byte oder */
#define EVENTSTORE_HEADER_VERSION	0x0001			/** Version number of Header file format */
#define EVENTSTORE_DATA_VERSION		0x0001			/** Version number of Data file format */

	/** A function which returns the last error code. */
#define ErrCode() 		errno
	/** An abstraction for a file handle. */
#define HANDLE 			int
	/** A value for an invalid file handle. */
#define INVALID_HANDLE_VALUE 						(-1)
	/** Get the size of a memory page for the system.
	 *	This is the basic size that the platform's memory manager uses, and is
	 *	fundamental to the use of memory-mapped files.
	 */
#define GET_PAGESIZE(x) 								((x) = sysconf(_SC_PAGE_SIZE))
#define EVENTSTORE_FLUSH_DATA_SYNC(file_descriptor) 	(!FlushFileBuffers(fd))
#define EVENTSTORE_MSYNC(addr, len, flags) 				(!FlushViewOfFile(addr, len))
#define EVENTSTORE_CLOSE_FILE(file_descriptor)			(CloseHandle(file_descriptor) ? 0 : -1)
#define EVENTSTORE_MEM_UNMAP(ptr, len)					UnmapViewOfFile(ptr)

#ifdef O_CLOEXEC /* Linux: Open file and set FD_CLOEXEC atomically */
#	define EVENTSTORE_CLOEXEC		O_CLOEXEC
#else
	 int fdflags;
#	define EVENTSTORE_CLOEXEC		0
#endif

#define READ_CREATE 			O_READ | O_CREAT | EVENTSTORE_CLOEXEC
#define READ_WRITE_CREATE		O_RDWR | O_CREAT | EVENTSTORE_CLOEXEC


#define HANDLE 							int 						/** An abstraction for a file handle. */
#define INVALID_HANDLE_VALUE 			(-1)						/** A value for an invalid file handle. */


/****************************************************************************
 *
 *		Helper functions
 *
 ****************************************************************************/


char debug_out[1000];
void vDebugPrint(const char* format, va_list args) {
	vsprintf(debug_out, format, args);
	DebugPrintf(debug_out);
}
void DebugPrint(const char* format, ...) {
	va_list args;
	va_start(args, format);
	vDebugPrint(format, args);
	va_end(args);
}


/****************************************************************************
 *
 *		Type Definitions
 *
 ****************************************************************************/
 
 typedef struct {
 	int 			error;			/** < Placeholder for error codes related to the file handle */
	char *			path;			/** < Physical path to the file referenced */
	HANDLE 			file_handle;	/** < Handle to the os file */
	char *			mmap_handle;	/** < Memory mapped view of the file */
	off_t 			file_size;		/** < Size of the file on disk */
	struct stat 	file_info;		/** < File info from the os */
} EVENTSTORE_file_handle;

/** Flag values for eventstore.flags field */
enum EVENTSTORE_flags
	{	F_OPEN_EXISTING			= 1 << 0
	,	F_CREATE_IF_MISSING		= 1 << 1
	,	F_READ_ONLY				= 1 << 2	/* Open database for read only access */
	,	F_NO_SYNC				= 1 << 3 	/* Don't fsync after commit */
	,	F_NO_SYNC_HEADER		= 1 << 4 	/* Don't fsync when writing to header metadata */
	};

typedef struct EVENTSTORE_settings EVENTSTORE_settings;
struct EVENTSTORE_settings {
	const char				header_version_label[8];
	const uint16_t			header_version;
	const char				data_version_label[8];
	const uint16_t			data_version;
};

	/** The EventStore database. */
typedef struct EVENTSTORE_database EVENTSTORE_database;
struct EVENTSTORE_database {
	pid_t			pid;				/** < Process ID of the database */
	uint32_t		flags;				/** < @ref Environment flags */
	int			 	page_size;			/** < Size of a page, from getpagesize(); */
	char * 			file_path;			/** < Path to the db files */

	EVENTSTORE_file_handle*			settings_file;		/** < Handle for database settings file */
	EVENTSTORE_file_handle*			header_file;		/** < Handle for managing the header file */
	EVENTSTORE_file_handle*			header_swap_file;	/** < Handle for mutable header swap file */
	EVENTSTORE_file_handle*			data_file;			/** < Handle for managing the data file */
	EVENTSTORE_file_handle*			generation_files;	/** < Handles for any active generation files */
	
	//EVENTSTORE_meta		*header_pages[2];	/** < Pointers to the two header pages */
	//pgno_t			max_pages;			/** < map_size / page_size */
	//EVENTSTORE_dbx		db_info;			/** < Array of static DB info */

	// Generations
	// Indexes
	// Memory Map(s)
	// Buffer(s)
};
// Disposal happens in close: mdb_env_close - Line: 4164



//HANDLE eventstore_open_file(char *path, int file_mode);



/****************************************************************************
 *
 *		Method Definitions
 *
 ****************************************************************************/
 
 void eventstore_write_file_info(char * handle) {
 	const char * header_version_label = "HEADER_V";
 	const char * data_version_label = "DATA___V";
 	int i = 0;
	handle[i] = (char)(EVENTSTORE_FILE_KEY >> 24); i++;
	handle[i] = (char)(EVENTSTORE_FILE_KEY >> 16); i++;
	handle[i] = (char)(EVENTSTORE_FILE_KEY >> 8); i++;
	handle[i] = (char)(EVENTSTORE_FILE_KEY); i++;
	int j = 0;
	int jsize = 0;

	jsize = sizeof(header_version_label);
	for(j = 0; j < jsize; j++) {
		handle[i] = header_version_label[j]; i++;
	}
	handle[i] = (char)(EVENTSTORE_HEADER_VERSION >> 8); i++;
	handle[i] = (char)(EVENTSTORE_HEADER_VERSION); i++;
	jsize = sizeof(data_version_label);
	for(j = 0; j < jsize; j++) {
		handle[i] = data_version_label[j]; i++;
	}
	handle[i] = (char)(EVENTSTORE_DATA_VERSION >> 8); i++;
	handle[i] = (char)(EVENTSTORE_DATA_VERSION); i++;
}

int open_file_handle(EVENTSTORE_file_handle* handle, char* dir_path, size_t dir_path_len, char* file_name, int file_name_len) {
	DebugPrint("File Size: %s", file_name);
	// Build the data file path by concat dir_path + file_name
	handle->path = (char*)malloc(dir_path_len + file_name_len);
	strcpy(handle->path, dir_path);
	strcat(handle->path, file_name);
	// TEMP: Delete the files first
	remove(handle->path);
	// Use the constructed path to try and open the db files
	handle->file_handle = open(handle->path, READ_WRITE_CREATE, 0666);
	// Make sure the file opened successfully
	if (handle->file_handle == INVALID_HANDLE_VALUE) {
		perror("error opening file");
		return 1;
	}
	// Grab the file's info
	if (fstat(handle->file_handle, &handle->file_info) == -1) {
		perror("error retrieving file info");
		return 1;
	}
	// Make sure it's actually a file
	if (!S_ISREG(handle->file_info.st_mode)) {
		fprintf (stderr, "%s is not a file\n", handle->path);
		return 1;
	}
	// Copy the size up to the structure (for easier access)
	handle->file_size = handle->file_info.st_size;
	DebugPrint("File Size: %d", handle->file_size);

	return 0;
}

EVENTSTORE_file_handle* open_file_handle2(char* dir_path, size_t dir_path_len, char* file_name, int file_name_len) {
	EVENTSTORE_file_handle* handle = malloc(sizeof(EVENTSTORE_file_handle));
	// Build the data file path by concat dir_path + file_name
	handle->path = (char*)malloc(dir_path_len + file_name_len);
	strcpy(handle->path, dir_path);
	strcat(handle->path, file_name);
	// TEMP: Delete the files first
	remove(handle->path);
	// Use the constructed path to try and open the db files
	handle->file_handle = open(handle->path, READ_WRITE_CREATE, 0666);
	// Make sure the file opened successfully
	if (handle->file_handle == INVALID_HANDLE_VALUE) {
		perror("error opening file");
		return handle;
	}
	// Grab the file's info
	if (fstat(handle->file_handle, &handle->file_info) == -1) {
		perror("error retrieving file info");
		return handle;
	}
	// Make sure it's actually a file
	if (!S_ISREG(handle->file_info.st_mode)) {
		fprintf (stderr, "%s is not a file\n", handle->path);
		return handle;
	}
	// Copy the size up to the structure (for easier access)
	handle->file_size = handle->file_info.st_size;
	DebugPrint("File Size: %d", handle->file_size);

	int init_file = 0;
	if (handle->file_size == 0) {
		init_file = 1;
		handle->file_size = 256;
		ftruncate(handle->file_handle, 256);
	}
	DebugPrint("File Size: %d", handle->file_size);
	int page_size = (int)getpagesize();
	DebugPrint("Page Size: %d", page_size);
	handle->mmap_handle = mmap(				// Create a memory map covering the file
		NULL, 							 	// Specifies the address to map into, -1 or NULL lets the system pick
		handle->file_size, 					// The size of the portion of the file being mapped to
		PROT_READ|PROT_WRITE,				// Access options PROT_READ | PROT_WRITE
		MAP_SHARED, 						// Shares map between processes
		handle->file_handle, 				// File handle to point the map at
		0);									// Offset from the start of the memory map
	if (handle->mmap_handle == (caddr_t)(-1)) {
		perror("error opening memory map");
		return handle;
	}
	DebugPrint("mmap_handle: %d", handle->mmap_handle);
	if (init_file > 0) {
		DebugPrint("Initialized: %s", handle->path);
		eventstore_write_file_info(handle->mmap_handle);
		//handle->mmap_handle[0] = (char)(0xC0);//DEC0DE';//(char)(EVENTSTORE_FILE_KEY & 0xFF000000);
		//handle->mmap_handle[1] = (char)(0xDE);
		//handle->mmap_handle[2] = (char)(EVENTSTORE_FILE_KEY & 0x0000FF00);
		//handle->mmap_handle[3] = (char)(EVENTSTORE_FILE_KEY & 0x000000FF);
		//memcpy(handle->mmap_handle, (char *)EVENTSTORE_FILE_KEY, sizeof(EVENTSTORE_FILE_KEY));
		//memset(handle->mmap_handle, EVENTSTORE_FILE_KEY, 32);
		//memset(handle->mmap_handle, EVENTSTORE_FILE_KEY, handle->file_size);
		//handle->mmap_handle[0] = '0';
	}
	DebugPrint("Byte at offset %d is [%c]\n", 2, handle->mmap_handle[0]);
	//uint64_t * long_mmap = handle->mmap_handle;
	//memset(handle->mmap_handle, 1, 2);

	return handle;
}


// DO_WRITE
// mdb_env_copyfd




void eventstore_open(char *path) { //, EVENTSTORE_flags *flags) {
	// Make the database struct instance
	EVENTSTORE_database* database = malloc(sizeof(EVENTSTORE_database));
	// Malloc a copy of the file handle structure for each db file
	database->settings_file = malloc(sizeof(EVENTSTORE_file_handle));
	database->header_file = malloc(sizeof(EVENTSTORE_file_handle));
	database->header_swap_file = malloc(sizeof(EVENTSTORE_file_handle));
	database->data_file = malloc(sizeof(EVENTSTORE_file_handle));
	// Get the size of the base path for str cat operations
	size_t path_size = strlen(path) + 1;
	// Populate file handles for all files
	int open_file_result = 0;
	open_file_result = open_file_handle(database->settings_file, path, path_size, SETTINGS_FILE_NAME, SETTINGS_FILE_NAME_SIZE);
	if (open_file_result > 0) { perror("Error opening settings file"); return; }
	open_file_result = open_file_handle(database->header_file, path, path_size, HEADER_FILE_NAME, HEADER_FILE_NAME_SIZE);
	if (open_file_result > 0) { perror("Error opening header file"); return; }
	open_file_result = open_file_handle(database->header_swap_file, path, path_size, HEADER_SWAP_FILE_NAME, HEADER_SWAP_FILE_NAME_SIZE);
	if (open_file_result > 0) { perror("Error opening header swap file"); return; }
	open_file_result = open_file_handle(database->data_file, path, path_size, DATA_FILE_NAME, DATA_FILE_NAME_SIZE);
	if (open_file_result > 0) { perror("Error opening data file"); return; }

	DebugPrint("Files opened");

	if (database->settings_file->file_size != sizeof(EVENTSTORE_settings)) {
		perror("Settings file is invalid");
		return;
	}

	errno = 0;
	return;

	//return database;
	//EVENTSTORE_file_handle* header_file = open_file_handle(path, path_size, HEADER_FILE_NAME, HEADER_FILE_NAME_SIZE);
	//EVENTSTORE_file_handle* data_file = open_file_handle(path, path_size, DATA_FILE_NAME, DATA_FILE_NAME_SIZE);

	//memcpy(NODEKEY(node), key->mv_data, key->mv_size);


	//int i = 0;
	//for (i = 0; i < 90; i++) {
		//write(header_file->file_handle, "stuff", 5);
		//memcpy(&data_file->mmap_handle, "qqqqq", 5);
		//header_file->mmap_handle[i] = i;
	//}
	/*char *header_file_path = (char*)malloc(path_size + 18);
	strcpy(header_file_path, path);
	strcat(header_file_path, "es_headerfile.esdb");
	HANDLE header_file = open(header_file_path, READ_WRITE_CREATE, 0666);
	// Unable to open file handle
	if (file_handle == INVALID_HANDLE_VALUE) {
		return INVALID_HANDLE_VALUE; // Line 3620 in mdb.c
	}
	// Seek from beginning to end to find file size
	off_t file_size = lseek(file_handle, 0, SEEK_END);
	if (file_size == -1) {
		return INVALID_HANDLE_VALUE;
	}
	void *header_map = mmap(NULL, )

	char *data_file_path = (char*)malloc(path_size + 16);
	strcpy(data_file_path, path);
	strcat(data_file_path, "es_datafile.esdb");
	HANDLE data_file = eventstore_open_file(data_file_path, 0666);

	free(header_file_path);
	free(data_file_path);*/

}





HANDLE eventstore_open_file(char *path, int file_mode) {
	HANDLE file_handle;
	off_t file_size;

	file_handle = open(path, READ_WRITE_CREATE, file_mode);

	// Unable to open file handle
	if (file_handle == INVALID_HANDLE_VALUE) {
		return INVALID_HANDLE_VALUE; // Line 3620 in mdb.c
	}

	return file_handle;
}


/*

int
mdb_put(MDB_txn *txn, MDB_dbi dbi,
    MDB_val *key, MDB_val *data, unsigned int flags)
{
	MDB_cursor mc;
	MDB_xcursor mx;

	assert(key != NULL);
	assert(data != NULL);

	if (txn == NULL || !dbi || dbi >= txn->mt_numdbs || !(txn->mt_dbflags[dbi] & DB_VALID))
		return EINVAL;

	if (F_ISSET(txn->mt_flags, MDB_TXN_RDONLY)) {
		return EACCES;
	}

	if (key->mv_size == 0 || key->mv_size > MDB_MAXKEYSIZE) {
		return EINVAL;
	}

	if ((flags & (MDB_NOOVERWRITE|MDB_NODUPDATA|MDB_RESERVE|MDB_APPEND|MDB_APPENDDUP)) != flags)
		return EINVAL;

	mdb_cursor_init(&mc, txn, dbi, &mx);
	return mdb_cursor_put(&mc, key, data, flags);
}
*/





