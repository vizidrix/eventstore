#include "eventstore.h"

// Uncomment this to wipe out files between runs
//#define ES_RESET_FILES

char *es_version(int *major, int *minor, int *patch)
{
	if (major) *major = ES_VERSION_MAJOR;
	if (minor) *minor = ES_VERSION_MINOR;
	if (patch) *patch = ES_VERSION_PATCH;
	return ES_VERSION_STRING;
}

#define ES_SETTINGS_FILE_NAME		"es_settings_file.esdb.txt"
#define ES_HEADER_FILE_NAME 		"es_header_file.esdb.txt"
#define ES_HEADER_SWAP_FILE_NAME	"es_header_swap_file.esdb.txt"
#define ES_DATA_FILE_NAME 			"es_data_file_%02d.esdb.txt"
#define ES_DATA_GEN_FILE_NAME		"es_data_gen_file_%02d.esdb.txt"

#define ES_SETTINGS_FILE_NAME_SIZE			sizeof(ES_SETTINGS_FILE_NAME) - 1
#define ES_HEADER_FILE_NAME_SIZE 			sizeof(ES_HEADER_FILE_NAME) - 1
#define ES_HEADER_SWAP_FILE_NAME_SIZE		sizeof(ES_HEADER_SWAP_FILE_NAME) - 1
#define ES_DATA_FILE_NAME_SIZE				sizeof(ES_DATA_FILE_NAME) - 1
#define ES_DATA_GEN_FILE_NAME_SIZE			sizeof(ES_DATA_GEN_FILE_NAME) - 1

/** Stamp to identify a file as ES valid and check for byte oder */
#define ES_FILE_KEY			"ES_HEADER_V_0x0001_DATA_V_0x0001"
#define ES_FILE_KEY_SIZE	sizeof(ES_FILE_KEY)

#define ES_LEVELS_IN_TRIE				 5	// Number of levels for trie depth
#define ES_MAX_GENERATIONS 				 4  // Number of generations allowed to be outstanding
// Addressing Trie:
// L5 ->   1 Gb *	64 blocks	=	 64 Gb (across 16 4Gb files) / DB + Header(s) + Generation(s) + Settings
// L4 ->  16 Mb	*	64		 	=	  1 Gb
// L3 -> 256 Kb	*	64		 	=	 16 Mb
// L2 ->   4 Kb	*	64		 	=	256 Kb
// L1 ->  64  b	*	64		 	=	  4 Kb
// --> Each address designates blocks of 64 bytes
// Addressed block = match L1 cache size
// --> File size = 1 Gb

// 6 bits = 0-63
// 5 levels @ 6 bits/lvl = 30 bits / 8 bits/byte = 4 bytes / address or uint32

/* Van Emde Boas layout
		        ____ ____ 1 ____ ____
		    **************************** Cut 2    
		      /                      \
		  __ 02 __                __ 03 __
		************************************* Cut 1
	    /          \            /          \
	   04          05          06          07
	 ****** C3   ****** C4   ****** C5   ****** C6  
	 /    \      /    \      /    \      /    \
	08    09    10 	  11    12    13    14    15

	01, 02, 03, 04, 08, 09, 05, 10, 11, 06, 12, 13, 07, 14, 15

*/
/*
	When copying up a shadow header each version must have a CRC Checksum.
	- Header data must be duplicated twice in the file at specified offsets.
	- Shadow header can be copied from either valid entries.
	- After shadow header is updated with a new Generation it is copied twice.
		- Copy over First Header.
		- Flush to disk.
		- Validate CRC Checksum on Frist Header.
		- Copy over Second Header.
		- Flush to disk.
		- Validate CRC Checksum on Second Header.
	* Should guarantee writes are persisted regardless of system interruption.
	See: http://guide.couchdb.org/draft/btree.html
*/

	/** Get the size of a memory page for the system.
	 *	This is the basic size that the platform's memory manager uses, and is
	 *	fundamental to the use of memory-mapped files.
	 */
#define GET_PAGESIZE(x) 								((x) = sysconf(_SC_PAGE_SIZE))
#define ES_FLUSH_DATA_SYNC(file_descriptor) 	(!FlushFileBuffers(fd))
#define ES_MSYNC(addr, len, flags) 				(!FlushViewOfFile(addr, len))
#define ES_CLOSE_FILE(file_descriptor)			(CloseHandle(file_descriptor) ? 0 : -1)
#define ES_MEM_UNMAP(ptr, len)					UnmapViewOfFile(ptr)

#ifdef O_CLOEXEC /* Linux: Open file and set FD_CLOEXEC atomically */
#	define ES_CLOEXEC		O_CLOEXEC
#else
	 int fdflags;
#	define ES_CLOEXEC		0
#endif

#define READ_CREATE 			O_READ | O_CREAT | ES_CLOEXEC
#define READ_WRITE_CREATE		O_RDWR | O_CREAT | ES_CLOEXEC


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
} ES_file_handle;

/** Flag values for eventstore.flags field */
enum ES_flags
	{	F_OPEN_EXISTING			= 1 << 0
	,	F_CREATE_IF_MISSING		= 1 << 1
	,	F_READ_ONLY				= 1 << 2	/* Open database for read only access */
	,	F_NO_SYNC				= 1 << 3 	/* Don't fsync after commit */
	,	F_NO_SYNC_HEADER		= 1 << 4 	/* Don't fsync when writing to header metadata */
	};

typedef struct ES_settings ES_settings;
struct ES_settings {
	unsigned char*		identifier[ES_FILE_KEY_SIZE];
 	//unsigned int			header_version;
	//unsigned int			data_version;
};

	/** The EventStore database. */
typedef struct ES_database ES_database;
struct ES_database {
	pid_t			pid;				/** < Process ID of the database */
	uint32_t		flags;				/** < @ref Environment flags */
	int			 	page_size;			/** < Size of a page, from getpagesize(); */
	char * 			file_path;			/** < Path to the db files */

	ES_file_handle*			settings_file;		/** < Handle for database settings file */
	ES_file_handle*			header_file;		/** < Handle for managing the header file */
	ES_file_handle*			header_swap_file;	/** < Handle for mutable header swap file */
	ES_file_handle*			data_files[16];		/** < Handle for managing the data file */
	ES_file_handle*			generation_files[ES_MAX_GENERATIONS];	/** < Handles for any active generation files */
	
	//ES_meta		*header_pages[2];	/** < Pointers to the two header pages */
	//pgno_t			max_pages;			/** < map_size / page_size */
	//ES_dbx		db_info;			/** < Array of static DB info */

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
 
 void es_write_settings_info(HANDLE handle) {
 	ftruncate(handle, ES_FILE_KEY_SIZE);
 	ssize_t count = write(handle, ES_FILE_KEY, ES_FILE_KEY_SIZE);

 	if (count != ES_FILE_KEY_SIZE) {
 		perror("Write to settings file failed");
 	}
}

int es_check_settings_info(HANDLE handle) {
	char data[ES_FILE_KEY_SIZE];
	lseek(handle, 0, SEEK_SET);
	read(handle, data, ES_FILE_KEY_SIZE);
	if (strcmp(data, ES_FILE_KEY) != 0) {
		perror("File key mismatch");
		return 1;
	}

	return 0;
}

int es_open_file_handle(ES_file_handle* handle, char* dir_path, size_t dir_path_len, char* file_name, int file_name_len) {
	// Build the data file path by concat dir_path + file_name
	handle->path = (char*)malloc(dir_path_len + file_name_len);
	strcpy(handle->path, dir_path);
	strcat(handle->path, file_name);
#ifdef ES_RESET_FILES
	// Delete the files first if the flag is set
	remove(handle->path);
#endif
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
	DebugPrint("Loaded file [%s] Size: %d Mb", handle->path, handle->file_size >> 20);

	return 0;
}

int es_verify_file_size(ES_file_handle* handle, off_t file_size) {
	if (handle->file_size == 0) {
		DebugPrint("Resizing to %d Mb", file_size >> 20);
		ftruncate(handle->file_handle, file_size);
		handle->file_size = file_size;
	}
	if (handle->file_size != file_size) {
		perror("file size out of sync");
		return 1;
	}

	return 0;
}

int es_open_mmap_handle(ES_file_handle* handle, off_t map_size) {
	handle->mmap_handle = mmap(				// Create a memory map covering the file
		NULL, 							 	// Specifies the address to map into, -1 or NULL lets the system pick
		map_size, 							// The size of the portion of the file being mapped to
		PROT_READ|PROT_WRITE,				// Access options PROT_READ | PROT_WRITE
		MAP_SHARED, 						// Shares map between processes
		handle->file_handle, 				// File handle to point the map at
		0);									// Offset from the start of the memory map
	if (handle->mmap_handle == (caddr_t)(-1)) {
		perror("error opening memory map");
		return 1;
	}

	return 0;
}

// DO_WRITE
// mdb_env_copyfd



void es_open(char *path) { //, ES_flags *flags) {
	// Make the database struct instance
	ES_database* database = malloc(sizeof(ES_database));
	
	// Get the size of the base path for str cat operations
	size_t path_size = strlen(path) + 1;
	// Populate file handles for all files
	int op_result = 0;

	// Malloc a copy of the file handle structure for each db file, open it, validate it, populate if needed

	database->settings_file = malloc(sizeof(ES_file_handle));
	op_result = es_open_file_handle(database->settings_file, path, path_size, ES_SETTINGS_FILE_NAME, ES_SETTINGS_FILE_NAME_SIZE);
	if (op_result > 0) { perror("Error opening settings file"); return; }
	// Validation checks for the settings file
	if (database->settings_file->file_size != sizeof(ES_settings)) {
		es_write_settings_info(database->settings_file->file_handle);
	}
	if (es_check_settings_info(database->settings_file->file_handle) > 0) {
		DebugPrint("Error in settings file");
		return;
	}
	DebugPrint("Settings file loaded");


	database->header_file = malloc(sizeof(ES_file_handle));
	op_result = es_open_file_handle(database->header_file, path, path_size, ES_HEADER_FILE_NAME, ES_HEADER_FILE_NAME_SIZE);
	if (op_result > 0) { perror("Error opening header file"); return; }
	op_result = es_verify_file_size(database->header_file, 1 << 20); // 1 MB
	if (op_result > 0) { perror("Error opening header file"); return; }
	//op_result = eventstore_open_mmap_handle(database->header_file, 1 << 20); // 1 MB
	//if (op_result > 0) { perror("Error opening header file"); return; }
	//database->header_swap_file = malloc(sizeof(ES_file_handle));
	DebugPrint("Header file loaded");

	int i = 0;
	char data_file_name[ES_DATA_FILE_NAME_SIZE];
	/*char* data_file_name = DATA_FILE_NAME;
	data_file_name[13] = '%';
	data_file_name[14] = 'd';*/
	for (i = 0; i < 16; i++) {
		database->data_files[i] = malloc(sizeof(ES_file_handle));

		sprintf(data_file_name, ES_DATA_FILE_NAME, i);
		//DebugPrint("Making data file: %s", data_file_name);
		op_result = es_open_file_handle(database->data_files[i], path, path_size, data_file_name, ES_DATA_FILE_NAME_SIZE);
		if (op_result > 0) { perror("Error opening data file"); return; }


		op_result = es_verify_file_size(database->data_files[i], ((off_t)1) << 32); // 1 GB * 4 = 2 ^ (30 + 2)
		if (op_result > 0) { perror("Error opening data file"); return; }
		//op_result = eventstore_open_mmap_handle(database->data_file, 1 << 36); // 1 GB * 32 = 2 ^ (30 + 6)
		//if (op_result > 0) { perror("Error opening data file"); return; }
	}
	DebugPrint("Data files loaded");

	// op_result = eventstore_open_file_handle(database->header_swap_file, path, path_size, HEADER_SWAP_FILE_NAME, HEADER_SWAP_FILE_NAME_SIZE);
	// if (op_result > 0) { perror("Error opening header swap file"); return; }

	errno = 0;
	return;

	//return database;
	//ES_file_handle* header_file = open_file_handle(path, path_size, HEADER_FILE_NAME, HEADER_FILE_NAME_SIZE);
	//ES_file_handle* data_file = open_file_handle(path, path_size, DATA_FILE_NAME, DATA_FILE_NAME_SIZE);

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






ES_file_handle* open_file_handle2(char* dir_path, size_t dir_path_len, char* file_name, int file_name_len) {
	ES_file_handle* handle = malloc(sizeof(ES_file_handle));
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
		//eventstore_write_settings_info(handle->mmap_handle);
		//handle->mmap_handle[0] = (char)(0xC0);//DEC0DE';//(char)(ES_FILE_KEY & 0xFF000000);
		//handle->mmap_handle[1] = (char)(0xDE);
		//handle->mmap_handle[2] = (char)(ES_FILE_KEY & 0x0000FF00);
		//handle->mmap_handle[3] = (char)(ES_FILE_KEY & 0x000000FF);
		//memcpy(handle->mmap_handle, (char *)ES_FILE_KEY, sizeof(ES_FILE_KEY));
		//memset(handle->mmap_handle, ES_FILE_KEY, 32);
		//memset(handle->mmap_handle, ES_FILE_KEY, handle->file_size);
		//handle->mmap_handle[0] = '0';
	}
	DebugPrint("Byte at offset %d is [%c]\n", 2, handle->mmap_handle[0]);

	return handle;
}


