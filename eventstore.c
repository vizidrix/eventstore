#include "fio.h"
#include "util.h"
#define __extern_golang
#include "../ringbuffer/ringbuffer.h"

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

#define MMAP_READ 						PROT_READ
#define MMAP_READ_WRITE					PROT_READ | PROT_WRITE

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

typedef struct ES_command {
	uint64_t		command_id;						/** <  First 56 bits are batch id, last 8 bits are command id */
	uint32_t		domain_id;						/** < Domain ID + Kind ID + Aggregate ID == 128 bits ~= UUID */
	uint32_t		kind_id;
	uint64_t		aggregate_id;
	uint32_t		crc;							/** < 32bit checksum of type+size+data */
	uint16_t		event_type;						/** < Identifies the structure in the data blob to client */
	uint16_t		event_size;						/** < Length of the event data */
	char			event_data[ES_MAX_DATA_SIZE];	/** < Bucket of data for the event */
} ES_command; // 8 + 4 + 4 + 8 + 4 + 2 + 2 = 32 bytes of overhead / command

typedef struct ES_data_barrier {
	uint64_t			seq_num;			/** < Index of the last entry released */
} ES_data_barrier;

typedef struct ES_write_buffer {
	uint64_t			seq_num;			/** < Index of next available entry in the buffer */
	ES_data_barrier *	producer_barrier;	
	ES_command *		command_buffer;		/** < Pre-allocated block of commands for batches */
	void *				data_buffer;		/** < Pre-allocated, contiguous block of data used to copy data in */
} ES_write_buffer;

// finished generation counter
// available generation counter <- seq_num of write buffer
// batch counter

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
	//writer->write_buffer = buffer;
	rb_buffer *buf;
	rb_buffer **buf_ptr = &buf;
	rb_init_buffer(buf_ptr, 2, 32); // New rules for sizing buffer
	//rb_init_buffer(buf_ptr, 16, 1024);
	rb_release_buffer(buf);
	DebugPrint("Made buffer: [ %d ]", sizeof(ES_command) * buffer_size);

}

void es_close_write_buffer(ES_writer* writer) {
	//free(writer->write_buffer->command_buffer);

	//free(writer->write_buffer);
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


struct ES_writer {
	char *					file_path;				/** < Path to the db files */

	fio_handle *			header_file;			/** < Handle for managing the header file */
	fio_handle *			data_files[16];			/** < Handles for managing the data file */
	
	ES_write_buffer *		write_buffer;
};

struct ES_batch_entry {
	uint16_t		* event_type;
	uint16_t		* event_size;
	char			* event_data[ES_MAX_DATA_SIZE];
};

struct ES_batch {
	uint64_t			batch_id;
	uint32_t			domain_id;
	uint32_t			kind_id;
	uint64_t			aggregate_id;
	char				buffer_size;
	char				batch_size;
	ES_batch_entry * 	entries;
};

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
	DebugPrint("Opening writer");
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


