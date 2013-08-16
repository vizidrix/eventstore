#include "eventstore.h"

/****************************************************************************
============================================================================================================
=        ==  ====  ==        ==  =======  ==        ===      ===        ====    ====       ===        ======
=  ========  ====  ==  ========   ======  =====  =====  ====  =====  ======  ==  ===  ====  ==  ============
=  ========  ====  ==  ========    =====  =====  =====  ====  =====  =====  ====  ==  ====  ==  ============
=  ========  ====  ==  ========  ==  ===  =====  ======  ==========  =====  ====  ==  ===   ==  ============
=      ====   ==   ==      ====  ===  ==  =====  ========  ========  =====  ====  ==      ====      ========
=  =========  ==  ===  ========  ====  =  =====  ==========  ======  =====  ====  ==  ====  ==  ============
=  =========  ==  ===  ========  =====    =====  =====  ====  =====  =====  ====  ==  ====  ==  ============
=  ==========    ====  ========  ======   =====  =====  ====  =====  ======  ==  ===  ====  ==  ============
=        =====  =====        ==  =======  =====  ======      ======  =======    ====  ====  ==        ======
============================================================================================================
****************************************************************************/

// Thanks for the ASCII comment blocks!  (Reverse font) http://patorjk.com/software/taag/

// Uncomment this to wipe out files between runs
//#define ES_RESET_FILES

char * es_version(int *major, int *minor, int *patch)
{
	if (major) *major = ES_VERSION_MAJOR;
	if (minor) *minor = ES_VERSION_MINOR;
	if (patch) *patch = ES_VERSION_PATCH;
	return ES_VERSION_STRING;
}

#define ES_SETTINGS_FILE_NAME		"es_settings_file.esdb.txt"
#define ES_HEADER_FILE_NAME 		"es_header_file.esdb.txt"
#define ES_DATA_GEN_FILE_NAME		"es_data_gen_file.esdb.txt"
#define ES_DATA_FILE_NAME 			"es_data_file_%02d.esdb.txt"


/** Stamp to identify a file as ES valid and check for byte oder */
#define ES_FILE_KEY			"ES_HEADER_V_0x0001_DATA_V_0x0001"
#define ES_FILE_KEY_SIZE	sizeof(ES_FILE_KEY)

#define FLAGS_READ_ONLY 				O_RDONLY
#define FLAGS_WRITE_ONLY				O_WRONLY
#define FLAGS_TRUNK						O_TRUNC
#define FLAGS_READ_CREATE 				O_READ | O_CREAT
#define FLAGS_READ_WRITE_CREATE			O_RDWR | O_CREAT

#define MODE_READ 						0444
#define MODE_WRITE 						0222
#define MODE_EXEC 						0111
#define MODE_READ_WRITE 				MODE_READ | MODE_WRITE
#define MODE_READ_WRITE_EXEC			MODE_READ | MODE_WRITE | MODE_EXEC

#define MMAP_READ 						PROT_READ
#define MMAP_READ_WRITE					PROT_READ|PROT_WRITE

#define HANDLE 							int 						/** An abstraction for a file handle. */
#define INVALID_HANDLE_VALUE 			(-1)

/****************************************************************************
=======================================================
=        ==  ====  ==       ===        ===      =======
====  =====   ==   ==  ====  ==  ========  ====  ======
====  ======  ==  ===  ====  ==  ========  ============
====  =======    ====       ===      =======  =========
====  ========  =====  ========  =============  =======
====  ========  =====  ========  ========  ====  ======
====  ========  =====  ========        ===      =======
=======================================================
****************************************************************************/



/****************************************************************************
=======================================================================
=       ===       ===    ==  ====  =====  =====        ==        ======
=  ====  ==  ====  ===  ===  ====  ====    =======  =====  ============
=  ====  ==  ===   ===  ===  ====  ==  ====  =====  =====  ============
=       ===      =====  ===   ==   ==  ====  =====  =====      ========
=  ========  ====  ===  ====  ==  ===        =====  =====  ============
=  ========  ====  ===  =====    ====  ====  =====  =====  ============
=  ========  ====  ==    =====  =====  ====  =====  =====        ======
=======================================================================
****************************************************************************/


/****************************************************************************
 *
 *		Database Types
 *
 ****************************************************************************/

	/** Abstraction which wraps both file handle and memmap handles */
typedef struct ES_file_handle {
 	int 			error;			/** < Placeholder for error codes related to the file handle */
	char *			path;			/** < Physical path to the file referenced */
	HANDLE 			file_handle;	/** < Handle to the os file */
	off_t 			file_size;		/** < Size of the file on disk */
	struct stat 	file_info;		/** < File info from the os */
} ES_file_handle;

typedef struct ES_mmap_handle {
	char *			mmap_handle;	/** < Memory mapped view of the file */
	off_t			size;			/** < Size of mmap overlay over file */
	off_t			offset;			/** < Byte offset from beginning of file */
} ES_mmap_handle;

//typedef struct ES_cursor {
//	char * 			dest;
//	off_t			size;
//	off_t			position;
//} ES_cursor;



typedef struct ES_put_event {
	char 			event_type; // Either [ BATCH_COMPLETE | COMMAND_COMPLETE ]
	uint64_t		id;
} ES_put_event;

/****************************************************************************
============================================================
=       ===  ====  ==      ===  ========    ====     =======
=  ====  ==  ====  ==  ===  ==  =========  ====  ===  ======
=  ====  ==  ====  ==  ===  ==  =========  ===  ============
=       ===  ====  ==      ===  =========  ===  ============
=  ========  ====  ==  ===  ==  =========  ===  ============
=  ========   ==   ==  ===  ==  =========  ====  ===  ======
=  =========      ===      ===        ==    ====     =======
============================================================
****************************************************************************/

struct ES_writer {
	char *					file_path;				/** < Path to the db files */

	ES_file_handle*			header_file;			/** < Handle for managing the header file */
	ES_file_handle*			generations_file;		/** < Handle for any active generations */
	ES_file_handle*			data_files[16];			/** < Handles for managing the data file */
	
	ES_mmap_handle *		generations_mmap;
};

struct ES_put_command {
	uint32_t		crc;							/** < 32bit checksum of type+data */
	uint64_t		command_id;						/** < First 56 bits are batch id, last 8 bits are command id */
	uint32_t		domain_id;						/** < Domain ID + Kind ID + Aggregate ID == 128 bits ~= UUID */
	uint32_t		kind_id;
	uint64_t		aggregate_id;
	uint16_t		event_type;						/** < Identifies the structure in the data blob to client */
	uint16_t		event_size;						/** < Length of the event data */
	char 			event_data[ES_MAX_DATA_SIZE];	/** < Bucket of data for the event */
}; // 4 + 8 + 4 + 4 + 8 + 2 + 2 = 32 bytes of overhead / command

struct ES_batch_entry {
	char			command_id;
	uint16_t		event_type;
	uint16_t		event_size;
	char *			event_data;
};

struct ES_batch {
	uint64_t			batch_id;
	uint32_t			domain_id;
	uint32_t			kind_id;
	uint64_t			aggregate_id;
	char				batch_size;
	ES_batch_entry * 	entries;
};

/****************************************************************************
============================================================================
=  =====  ==        ==        ==  ====  ====    ====       ====      =======
=   ===   ==  ===========  =====  ====  ===  ==  ===  ====  ==  ====  ======
=  =   =  ==  ===========  =====  ====  ==  ====  ==  ====  ===  ===========
=  == ==  ==      =======  =====        ==  ====  ==  ====  =====  =========
=  =====  ==  ===========  =====  ====  ==  ====  ==  ====  =======  =======
=  =====  ==  ===========  =====  ====  ===  ==  ===  ====  ==  ====  ======
=  =====  ==        =====  =====  ====  ====    ====       ====      =======
============================================================================
****************************************************************************/



/****************************************************************************
=======================================================================
=       ===       ===    ==  ====  =====  =====        ==        ======
=  ====  ==  ====  ===  ===  ====  ====    =======  =====  ============
=  ====  ==  ===   ===  ===  ====  ==  ====  =====  =====  ============
=       ===      =====  ===   ==   ==  ====  =====  =====      ========
=  ========  ====  ===  ====  ==  ===        =====  =====  ============
=  ========  ====  ===  =====    ====  ====  =====  =====  ============
=  ========  ====  ==    =====  =====  ====  =====  =====        ======
=======================================================================
****************************************************************************/



/****************************************************************************
 *
 *		Helper Methods
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
 *		Disk and Memory Methods
 *
 ****************************************************************************/

int es_open_file_handle(ES_file_handle** handle_ptr, char* dir_path, char* file_name, int flags, int mode) {
	// Allocate the memory to hold the handle
	*handle_ptr = malloc(sizeof(ES_file_handle));
	ES_file_handle* handle = *handle_ptr;
	// Get the size of the base path and file for str cat operations
	size_t dir_path_len = strlen(dir_path) + 1;
	size_t file_name_len = strlen(file_name) + 1;
	// Build the data file path by concat dir_path + file_name
	handle->path = (char*)malloc(dir_path_len + file_name_len);
	strcpy(handle->path, dir_path);
	strcat(handle->path, file_name);
#ifdef ES_RESET_FILES
	// Delete the files first if the flag is set
	remove(handle->path);
#endif
	// Use the constructed path to try and open the db files
	handle->file_handle = open(handle->path, flags, mode);
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
	//DebugPrint("Loaded file [%s] Size: %d Mb", handle->path, handle->file_size >> 20);

	return 0;
}

void es_close_file_handle(ES_file_handle* handle) {
	// Make sure file reference is closed
	close(handle->file_handle);
	// Free the path string
	free(handle->path);
	// Free the handle itself
	free(handle);
}

int es_verify_file_size(ES_file_handle* handle, off_t file_size) {
	if (handle->file_size < file_size) {
		ftruncate(handle->file_handle, file_size);
		handle->file_size = file_size;
	}
	if (handle->file_size != file_size) {
		perror("file size out of sync");
		return 1;
	}

	return 0;
}

int es_open_mmap_handle(ES_mmap_handle** handle_ptr, ES_file_handle* file, off_t offset, off_t size, int mmap_mode) {
	// Allocate the memory to hold the handle
	*handle_ptr = malloc(sizeof(ES_mmap_handle));
	ES_mmap_handle* handle = *handle_ptr;
	
	handle->mmap_handle = mmap(				// Create a memory map covering the file
		NULL, 							 	// Specifies the address to map into, -1 or NULL lets the system pick
		size, 								// The size of the portion of the file being mapped to
		mmap_mode,							// Access options PROT_READ | PROT_WRITE
		MAP_SHARED, 						// Shares map between processes
		file->file_handle, 					// File handle to point the map at
		offset);							// Offset from the start of the memory map
	if (handle->mmap_handle == (caddr_t)(-1)) {
		perror("error opening memory map");
		return 1;
	}
	return 0;
}

void es_close_mmap_handle(ES_mmap_handle* handle, off_t size) {
	// Release memory held by mmap
	munmap(handle->mmap_handle, size);
	// Free the handle itself
	free(handle);
}

/****************************************************************************
============================================================
=       ===  ====  ==      ===  ========    ====     =======
=  ====  ==  ====  ==  ===  ==  =========  ====  ===  ======
=  ====  ==  ====  ==  ===  ==  =========  ===  ============
=       ===  ====  ==      ===  =========  ===  ============
=  ========  ====  ==  ===  ==  =========  ===  ============
=  ========   ==   ==  ===  ==  =========  ====  ===  ======
=  =========      ===      ===        ==    ====     =======
============================================================
****************************************************************************/


/****************************************************************************
 *
 *		Writer Methods
 *
 ****************************************************************************/

ES_writer* es_open_write(char* path) {
	DebugPrint("Opening writer");
	int op_result = 0;
	// Allocate memory for use by the writer struct
	ES_writer * writer = malloc(sizeof(ES_writer));
	// Open the file in Read/Write/Create mode
	op_result = es_open_file_handle(&writer->generations_file, path, ES_DATA_GEN_FILE_NAME, FLAGS_READ_WRITE_CREATE, MODE_READ_WRITE);
	if (op_result > 0) { perror("Error opening generations file"); return NULL; }


	int generations_file_header_size = 1024;
	int generations_file_data_size = 65536;//4096;
	int max_generations = 6;
	int generations_file_size = generations_file_header_size + (generations_file_data_size * max_generations);
	es_verify_file_size(writer->generations_file, generations_file_size);
	
	op_result = es_open_mmap_handle(&writer->generations_mmap, writer->generations_file, 0, generations_file_size, MMAP_READ_WRITE);
	if (op_result > 0) { perror("Error opening generations map"); return NULL; }

	// finish opening the mmap
	// use the mmap data as a command source
	// populate it using disruptor style allocator
	// integrate with Go mapped slices

	//writer->commands = malloc(sizeof(ES_put_command) * 8);

	return writer;
}

void es_close_write(ES_writer* writer) {
	//free(writer->commands);
	es_close_mmap_handle(writer->generations_mmap, writer->generations_file->file_size);
	es_close_file_handle(writer->generations_file);
	free(writer);
	DebugPrint("Closed writer");
}

void es_publish_batch(ES_batch* batch) {
	// perform publish actions
	int i = 0;
	for(i = 0; i < batch->batch_size; i++) {
		free(batch->entries[i].event_data);
	}
	free(batch->entries);
	free(batch);
}

ES_batch* es_alloc_batch(ES_writer* writer, 
	uint32_t domain_id, 
	uint32_t kind_id, 
	uint64_t aggregate_id, 
	char count) {
	DebugPrint("[es.....c]\tAllocating Batch: %d", count);

	// Need to make a batch
	ES_batch* batch = malloc(sizeof(ES_batch));
	batch->batch_id = 5;
	batch->domain_id = domain_id;
	batch->kind_id = kind_id;
	batch->aggregate_id = aggregate_id;
	batch->batch_size = count;
	batch->entries = malloc(sizeof(ES_batch_entry) * count);

	int i, j = 0;
	for(i = 0; i < count; i++) {
		batch->entries[i].command_id = i;
		batch->entries[i].event_type = 1;
		batch->entries[i].event_size = 0;
		// Need to change this to point to the correct slot in generation mmap
		batch->entries[i].event_data = malloc(10);//ES_MAX_DATA_SIZE);
		//char * data = &batch->entries[i].event_data
		//char * data = malloc(10);//ES_MAX_DATA_SIZE);
		//char * data = malloc(10);
		for(j = 0; j < 10; j++) {
			batch->entries[i].event_data[j] = j * i;
		}
		//&batch->entries[i].event_data = &data;
	}
	// Populate batch with number of entries

	return batch;
	//DebugPrint("Size: %d", sizeof(ES_put_command));

	//ES_put_command* commands = (ES_put_command*)writer->generations_mmap->mmap_handle;

	//return commands;
/*

	//void* mem = (struct S_put_command*)malloc(sizeof(ES_put_command) * 4);
	//struct ES_put_command *commands[2];// = malloc(sizeof(ES_put_command) * 2);// = ES_put_command[1];
	//commands[0] = malloc(sizeof(ES_put_command));
	writer->commands[0].batch_id = 10;
	writer->commands[0].command_id = 11;
	writer->commands[0].crc = 12;
	writer->commands[0].event_data[0] = 99;
	writer->commands[0].event_data[20] = 44;

	return writer->commands;

	*/

/*
	
	ES_put_command * commands;
	commands = malloc(sizeof(ES_put_command) * 2);
	//commands[0] = &mem[0];
	commands[0].batch_id = 10;
	commands[0].command_id = 11;
	commands[0].crc = 12;
	commands[0].event_data[0] = 44;
	//commands[0].event_data[1] = 43;
	//commands[1] = malloc(sizeof(ES_put_command));
	//commands[1] = &mem[64];
	commands[1].batch_id = 20;
	commands[1].command_id = 21;

	ES_put_command * temp = &commands[0];
	return temp;
	//return commands;//[0];//[0];
	*/
	
}













