#include "fio.h"
#include "eventstore.h"

/***********************************************************************************************************
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
***********************************************************************************************************/

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
#define ES_FILE_KEY						"ES_HEADER_V_0x0001_DATA_V_0x0001"
#define ES_FILE_KEY_SIZE				sizeof(ES_FILE_KEY)

#define MMAP_READ 						PROT_READ
#define MMAP_READ_WRITE					PROT_READ | PROT_WRITE

/******************************************************
=======================================================
=        ==  ====  ==       ===        ===      =======
====  =====   ==   ==  ====  ==  ========  ====  ======
====  ======  ==  ===  ====  ==  ========  ============
====  =======    ====       ===      =======  =========
====  ========  =====  ========  =============  =======
====  ========  =====  ========  ========  ====  ======
====  ========  =====  ========        ===      =======
=======================================================
******************************************************/



/**********************************************************************
=======================================================================
=       ===       ===    ==  ====  =====  =====        ==        ======
=  ====  ==  ====  ===  ===  ====  ====    =======  =====  ============
=  ====  ==  ===   ===  ===  ====  ==  ====  =====  =====  ============
=       ===      =====  ===   ==   ==  ====  =====  =====      ========
=  ========  ====  ===  ====  ==  ===        =====  =====  ============
=  ========  ====  ===  =====    ====  ====  =====  =====  ============
=  ========  ====  ==    =====  =====  ====  =====  =====        ======
=======================================================================
**********************************************************************/


/****************************************************************************
 *
 *		Database Types
 *
 ****************************************************************************/

typedef struct ES_data_barrier {
	uint64_t			seq_num;			/** < Index of the last entry released */
} ES_data_barrier;

typedef struct ES_write_buffer {
	uint64_t			seq_num;			/** < Index of next available entry in the buffer */
	ES_data_barrier *	producer_barrier;	
	ES_command *		command_buffer;		/** < Pre-allocated block of commands for batches */
	void *				data_buffer;		/** < Pre-allocated, contiguous block of data used to copy data in */
} ES_write_buffer;


// Write to ring -> 
//		calc CRC -> 
//			scan forward until
//			a) as far as possible with certain batch end reached
//				(must see next batch id to know current is finished)
//			b) memory buffer exceeds 4k buffer (roll back to last batch id)
//			- Write batch off to disk through data distribution logic
//			- Identify this batch as a "generation" to enable generational read/index concept
//		update index

//typedef struct ES_mmap_handle {
//	void *			mmap_handle;	/** < Memory mapped view of the file */
//	off_t			size;			/** < Size of mmap overlay over file */
//	off_t			offset;			/** < Byte offset from beginning of file */
//} ES_mmap_handle;



//typedef struct ES_put_event {
//	char 			event_type; // Either [ BATCH_COMPLETE | COMMAND_COMPLETE ]
//	uint64_t		id;
//} ES_put_event;


	/** Root management structure for the gen file */
//typedef struct ES_gen_file {
//	uint64_t		done_gen_index;			/** < Index of the last completed gen */
//	uint64_t		next_gen_index;			/** < Index to be allocated to the next gen */
//	uint64_t		batch_counter;			/** < Batches are the lower 24 bits shifted by 8 to make room for batch count */

//} ES_gen_file;

/*
New id's are seperated from existing id's
- New id's can be streamed in as a batch
- Existing id's can be split between quick append and realloc append
	- If there is room in the block then just append to existing data
	- If a data move is required then
		- Do the relocation(s) in batche(s)
		- Do the append into the newly available space


On any publish, if data size of batch will exceed targeted
block size (4k) for new id set then clamp the current gen 
and start writing to the next

http://www.cse.ohio-state.edu/~zhang/hpca11-submitted.pdf

*/
//typedef struct ES_gen {
//	uint64_t		gen_index;			/** Index assigned to this page */
//	char			page_count;			/** Number of pages in this gen */
//
//	uint16_t		data_size;			/** Cumulative raw data size in the insert of this gen */
//
//} ES_gen;

	/** Page management structure for the gen file */
//typedef struct ES_gen_page {
//
//	uint64_t		event_count;		/** Number of events in the page */
//	void *			mmap_handle;		/** Pointer to this section of the generation file */
//} ES_gen_page;

//typedef struct ES_cursor {
//	char * 			dest;
//	off_t			size;
//	off_t			position;
//} ES_cursor;



/***********************************************************
============================================================
=       ===  ====  ==      ===  ========    ====     =======
=  ====  ==  ====  ==  ===  ==  =========  ====  ===  ======
=  ====  ==  ====  ==  ===  ==  =========  ===  ============
=       ===  ====  ==      ===  =========  ===  ============
=  ========  ====  ==  ===  ==  =========  ===  ============
=  ========   ==   ==  ===  ==  =========  ====  ===  ======
=  =========      ===      ===        ==    ====     =======
============================================================
***********************************************************/

struct ES_writer {
	char *					file_path;				/** < Path to the db files */

	//uint64_t				next_gen_counter;		/** < Counter is the next generation id to issue */
	//uint64_t				completed_gen_counter;	/** < The last completed generation id */
	//uint64_t				batch_counter;			/** < Batches are the lower 24 bits shifted by 8 to make room for batch count */

	fio_handle*			header_file;			/** < Handle for managing the header file */
	//ES_file_handle*			header_file;			/** < Handle for managing the header file */
	//ES_file_handle*			gen_file;				/** < Handle for any active generations */
	//ES_file_handle*			data_files[16];			/** < Handles for managing the data file */
	fio_handle*			data_files[16];			/** < Handles for managing the data file */
	
	//ES_mmap_handle *		gen_mmap;				/** < Mmap over the header of the gen file */

	ES_write_buffer *		write_buffer;
};

struct ES_command {
	uint64_t		command_id;						/** <  First 56 bits are batch id, last 8 bits are command id */
	uint32_t		domain_id;						/** < Domain ID + Kind ID + Aggregate ID == 128 bits ~= UUID */
	uint32_t		kind_id;
	uint64_t		aggregate_id;
	uint32_t		crc;							/** < 32bit checksum of type+size+data */
	uint16_t		event_type;						/** < Identifies the structure in the data blob to client */
	uint16_t		event_size;						/** < Length of the event data */
	char			event_data[ES_MAX_DATA_SIZE];	/** < Bucket of data for the event */
};

//typedef struct ES_put_command {
//	uint32_t		crc;							/** < 32bit checksum of type+data */
//	uint64_t		command_id;						/** < First 56 bits are batch id, last 8 bits are command id */
//	uint32_t		domain_id;						/** < Domain ID + Kind ID + Aggregate ID == 128 bits ~= UUID */
//	uint32_t		kind_id;
//	uint64_t		aggregate_id;
//	uint16_t		event_type;						/** < Identifies the structure in the data blob to client */
//	uint16_t		event_size;						/** < Length of the event data */
//	char 			event_data[ES_MAX_DATA_SIZE];	/** < Bucket of data for the event */
//} ES_put_command; // 4 + 8 + 4 + 4 + 8 + 2 + 2 = 32 bytes of overhead / command

//struct ES_batch_entry {
//	char			command_id;
//	uint16_t		event_type;
//	uint16_t		event_size;
//	uint32_t		crc;
//	off_t			mmap_offset;
//	char *			event_data;
//};

//struct ES_batch {
//	uint64_t			batch_id;
//	uint32_t			domain_id;
//	uint32_t			kind_id;
//	uint64_t			aggregate_id;
//	char				buffer_size;
//	char				batch_size;
//	ES_batch_entry * 	entries;
//};

/***************************************************************************
============================================================================
=  =====  ==        ==        ==  ====  ====    ====       ====      =======
=   ===   ==  ===========  =====  ====  ===  ==  ===  ====  ==  ====  ======
=  =   =  ==  ===========  =====  ====  ==  ====  ==  ====  ===  ===========
=  == ==  ==      =======  =====        ==  ====  ==  ====  =====  =========
=  =====  ==  ===========  =====  ====  ==  ====  ==  ====  =======  =======
=  =====  ==  ===========  =====  ====  ===  ==  ===  ====  ==  ====  ======
=  =====  ==        =====  =====  ====  ====    ====       ====      =======
============================================================================
***************************************************************************/



/**********************************************************************
=======================================================================
=       ===       ===    ==  ====  =====  =====        ==        ======
=  ====  ==  ====  ===  ===  ====  ====    =======  =====  ============
=  ====  ==  ===   ===  ===  ====  ==  ====  =====  =====  ============
=       ===      =====  ===   ==   ==  ====  =====  =====      ========
=  ========  ====  ===  ====  ==  ===        =====  =====  ============
=  ========  ====  ===  =====    ====  ====  =====  =====  ============
=  ========  ====  ==    =====  =====  ====  =====  =====        ======
=======================================================================
**********************************************************************/



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
 *		Memory Methods
 *
 ****************************************************************************/

void es_init_write_buffer(ES_writer* writer, uint32_t buffer_size) {
	ES_write_buffer * buffer = malloc(sizeof(ES_write_buffer));

	buffer->seq_num = 0;
	// Allocate enough room for the specified number of entries
	buffer->command_buffer = malloc(sizeof(ES_command) * buffer_size);
	//buffer->data_buffer = malloc(data_buffer_size);
	writer->write_buffer = buffer;
	DebugPrint("Made buffer: [ %d ]", sizeof(ES_command) * buffer_size);

}

void es_close_write_buffer(ES_writer* writer) {
	free(writer->write_buffer->command_buffer);

	free(writer->write_buffer);
}
/*
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
*/

//uint64_t es_get_batch_id(ES_writer* writer) {
//	writer->batch_counter++;
//	return writer->batch_counter;
//}

/*
void es_print_batch(ES_batch* batch) {
	DebugPrint("** Publishing **");
	DebugPrint("Batch Id: %d", batch->batch_id);
	DebugPrint("Domain Id: %d", batch->domain_id);
	DebugPrint("Kind Id: %d", batch->kind_id);
	DebugPrint("Aggregate Id: %d", batch->aggregate_id);
	DebugPrint("Batch Size: %d", batch->batch_size);

	int i = 0;
	for(i = 0; i < batch->batch_size; i++) {
		DebugPrint("\tCommand Id: %d", batch->entries[i].command_id);
		DebugPrint("Event Type: %d", batch->entries[i].event_type);
		DebugPrint("Event Size: %d", batch->entries[i].event_size);
	}
}
*/

void es_load_gen_header(ES_writer* writer) {
	int op_result = 0;

	int generations_file_header_size = 1024;
	int generations_file_data_size = 65536;//4096;
	int max_generations = 6;
	//int generations_file_size = generations_file_header_size + (generations_file_data_size * max_generations);
	//es_verify_file_size(writer->gen_file, generations_file_size);
	
	// Load the batch_id from storage on load
	//void* gen_map =mmap(NULL, 4096, MMAP_READ_WRITE, MAP_SHARED, writer->gen_file->file_handle, 0);
	
	//void* gen_map = mmap(NULL, 4096, MMAP_READ_WRITE, MAP_SHARED, writer->gen_file->file_handle, 0);
	//ES_gen_file *gen = (ES_gen_file*)gen_map;
	//gen->done_gen_index = 300;
	//gen->next_gen_index = 400;
	//DebugPrint("Gen: %d", gen->done_gen_index);
	//DebugPrint("Gen: %d", gen->next_gen_index);


	//op_result = es_open_mmap_handle(
	//	&writer->gen_mmap,				/** The mmap handle */
	//	writer->gen_file,				/** File handle to point at */
	//	0,								/** Offset within the file to start at */
	//	4096,//sizeof(ES_gen_file)*2, //generations_file_header_size,	/** Size of the section to map over */
	//	MMAP_READ_WRITE);				/** Enable both read and write to the map */
	//if (op_result > 0) { perror("Error opening generations map"); return; }

	/*
	ES_gen_file *gen = (ES_gen_file*)writer->gen_mmap->mmap_handle;
	gen->done_gen_index = 1000;
	gen->next_gen_index = 22;
	//DebugPrint("Gen: %s", writer->gen_mmap->mmap_handle[0]);
	DebugPrint("Gen: %d", gen->done_gen_index);
	DebugPrint("Gen: %d", gen->next_gen_index);
	*/
}

/***********************************************************
============================================================
=       ===  ====  ==      ===  ========    ====     =======
=  ====  ==  ====  ==  ===  ==  =========  ====  ===  ======
=  ====  ==  ====  ==  ===  ==  =========  ===  ============
=       ===  ====  ==      ===  =========  ===  ============
=  ========  ====  ==  ===  ==  =========  ===  ============
=  ========   ==   ==  ===  ==  =========  ====  ===  ======
=  =========      ===      ===        ==    ====     =======
============================================================
***********************************************************/


/****************************************************************************
 *
 *		Writer Methods
 *
 ****************************************************************************/


// 2 Mb / generation * N generations
// 64 pages per generation = 512 events / generation
// 4k page / 512 = 8 events / page
// 4096 Mb -> 4 Gb gen file + header
// 8388608 slots (32768 * 256 max batch size) @ 512 bytes

ES_writer* es_open_writer(char* path) {
	//DebugPrint("Opening writer");
	int op_result = 0;
	// Allocate memory for use by the writer struct
	ES_writer * writer = malloc(sizeof(ES_writer));
	// Open the file in Read/Write/Create mode
	//op_result = es_open_file_handle(&writer->gen_file, path, ES_DATA_GEN_FILE_NAME, FLAGS_READ_WRITE_CREATE, MODE_READ_WRITE);
	//if (op_result > 0) { perror("Error opening generations file"); return NULL; }


	es_init_write_buffer(writer, 1<<12/*4096*/);//, 1<<12<<12/*4096*4096*/);

	//es_load_gen_header(writer);

	// finish opening the mmap
	// use the mmap data as a command source
	// populate it using disruptor style allocator
	// integrate with Go mapped slices

	//writer->commands = malloc(sizeof(ES_put_command) * 8);

	return writer;
}

void es_close_writer(struct ES_writer* writer) {
	es_close_write_buffer(writer);
	//free(writer->commands);
	//es_close_mmap_handle(writer->generations_mmap, writer->generations_file->file_size);
	//es_close_file_handle(writer->gen_file);
	free(writer);
	//DebugPrint("Closed writer");
}

ES_batch* es_alloc_batch(struct ES_writer* writer, 
	uint32_t domain_id, 
	uint32_t kind_id, 
	uint64_t aggregate_id,
	char count) {
	// Get batch id

	// Alloc count commands from buffer
	// Populate default values
	// Empty bucket
}
/*
ES_batch* es_alloc_batch(ES_writer* writer, 
	uint32_t domain_id, 
	uint32_t kind_id, 
	uint64_t aggregate_id,
	char size,
	char count) {
	DebugPrint("[es.....c]\tAllocating Batchs: %d of %d", count, 1<<size);
	if(size > 12) {
		perror("Max buffer size is 4096");
		return NULL;
	}

	// Need to make a batch
	ES_batch* batch = malloc(sizeof(ES_batch));
	batch->batch_id = 0;//es_get_batch_id(writer);
	batch->domain_id = domain_id;
	batch->kind_id = kind_id;
	batch->aggregate_id = aggregate_id;
	batch->buffer_size = size;
	batch->batch_size = count;
	batch->entries = malloc(sizeof(ES_batch_entry) * count);

	int i = 0;
	for(i = 0; i < count; i++) {
		batch->entries[i].command_id = i;
		batch->entries[i].event_type = 0;
		batch->entries[i].event_size = 0;
		// Need to change this to point to the correct slot in generation mmap
		batch->entries[i].event_data = malloc(1 << size);

		//for(j = 0; j < ES_MAX_DATA_SIZE; j++) {
		//	batch->entries[i].event_data[j] = j * i;
		//}
	}
	// Populate batch with number of entries

	return batch;
}
*/

void es_publish_batch(ES_batch* batch) {
	// perform publish actions
	//es_print_batch(batch);
	//int i = 0;
	//for(i = 0; i < batch->batch_size; i++) {
	//	free(batch->entries[i].event_data);
	//}
	//free(batch->entries);
	//free(batch);
	
}


