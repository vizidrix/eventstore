#ifndef _FIO_H_
#define _FIO_H_

#include <stdarg.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#ifdef HAVE_SYS_FILE_H
#include <sys/file.h>
#endif
#include <fcntl.h>
#include <sys/stat.h>

#define FIO_SUCCESS 					0 						/** Successful result */
#define FIO_ERROR						(-30600)				/** Generic error */

#define FIO_NOTFOUND 					(FILEIO_ERROR - 100)	/** File was not found */
#define FIO_INVALID 					(FILEIO_ERROR - 200)	/** File is invalid */

#define FIO_HANDLE						int						/** An abstraction for a file handle. */
#define FIO_INVALID_HANDLE				(-1)

#define FIO_FLAGS_READ_ONLY 			O_RDONLY
#define FIO_FLAGS_WRITE_ONLY			O_WRONLY
#define FIO_FLAGS_TRUNK					O_TRUNC
#define FIO_FLAGS_READ_CREATE 			O_READ | O_CREAT
#define FIO_FLAGS_READ_WRITE_CREATE		O_RDWR | O_CREAT

#define FIO_MODE_READ					0444
#define FIO_MODE_WRITE					0222
#define FIO_MODE_EXEC					0111
#define	FIO_MODE_READ_WRITE				FILEIO_READ | FILEIO_WRITE
#define FIO_MODE_READ_WRITE_EXEC		FILEIO_READ | FILEIO_WRITE | FILEIO_EXEC

/****************************************************************************
 *
 *		Signatures
 *
 ****************************************************************************/

/** Structure which wraps file handle */
typedef struct fio_handle fio_handle;

int fio_open(fio_handle** handle_ptr, const char* dir_path, const char* file_name, int flags, int mode);
int fio_close(fio_handle* handle);
int fio_set_size(fio_handle* handle, off_t file_size);

/****************************************************************************
 *
 *		Implementations
 *
 ****************************************************************************/

struct fio_handle {
 	int 			error;			/** < Placeholder for error codes related to the file handle */
	char *			path;			/** < Physical path to the file referenced */
	FIO_HANDLE 		file_handle;	/** < Handle to the os file */
	off_t 			file_size;		/** < Size of the file on disk */
};

int fio_open(struct fio_handle** handle_ptr, const char* dir_path, const char* file_name, int flags, int mode) {
	// Allocate the memory to hold the handle
	*handle_ptr = malloc(sizeof(fio_handle));
	fio_handle* handle = *handle_ptr;
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
	if (handle->file_handle == FIO_INVALID_HANDLE) {
		perror("error opening file");
		return 1;
	}
	// Grab the file's info
	struct stat file_info;
	if (fstat(handle->file_handle, &file_info) == -1) {
		perror("error retrieving file info");
		return 1;
	}
	// Make sure it's actually a file
	if (!S_ISREG(file_info.st_mode)) {
		fprintf (stderr, "%s is not a file\n", handle->path);
		return 1;
	}
	// Copy the size up to the structure (for easier access)
	handle->file_size = file_info.st_size;
	//DebugPrint("Loaded file [%s] Size: %d Mb", handle->path, handle->file_size >> 20);

	return 0;
}

int fio_close(struct fio_handle* handle) {
	// Make sure file reference is closed
	close(handle->file_handle);
	// Free the path string
	free(handle->path);
	// Free the handle itself
	free(handle);

	return 0;
}

int fio_set_size(struct fio_handle* handle, off_t file_size) {
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

#endif