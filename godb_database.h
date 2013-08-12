#ifndef godb_database_h
#define godb_database_h

#include "godb.h"
//#include <sys/types.h>
//#include <stdint.h>

//#include <sys/types.h>
//#include <string.h>

//#include <stdio.h>

//#include <sys/stat.h>
//#include <sys/param.h>

//#include <sys/uio.h>
//#include <sys/mman.h>
//#ifdef HAVE_SYS_FILE_H
#include <sys/file.h>
//#endif
//#include <fcntl.h>

//#include <assert.h>
//#include <errno.h>
//#include <limits.h>
//#include <stddef.h>
//#include <inttypes.h>
//#include <stdio.h>

//#include <string.h>
////#include <time.h>
//#include <unistd.h>

#define GODB_FILE_KEY			0xC0DEC0DE		/** Stamp to identify a file as GODB valid and check for byte oder */
#define GODB_HEADER_VERSION		0x0001			/** Version number of Header file format */
#define GODB_DATA_VERSION		0x0001			/** Version number of Data file format */


//#define GODB_WRITEMAP  0x0001* Use writable mmap */
//#define GODB_MAPASYNC 0x0001* Use async msync when GODB_WRITEMAP is used */				

/** @ defgroup errors 	Return Codes
*
*	Avoid conflict with BerkelyDB (-30800 to -30999) and MDB (-30799 to -30783)
*	@{
*/
#define GODB_SUCCESS 			0 						/** Successful result */
#define GODB_ERROR				(-30600)				/** Generic error */
#define GODB_NOTFOUND 			(GODB_ERROR - 1) 		/** Key not found during get (EOF) */
#define GODB_CORRUPTED 			(GODB_ERROR - 2)		/** CheckSum failed */
#define GODB_PANIC 				(GODB_ERROR - 3) 		/** Update of meta page failed, probably I/O error */
#define GODB_VERSION_MISMATCH 	(GODB_ERROR - 4) 		/** Environment version mismatch */
#define GODB_INVALID 			(GODB_ERROR - 5)		/** Invalid GODB file */
#define GODB_MAP_FULL			(GODB_ERROR - 6)		/** Environment mapsize reached */
#define GODB_PAGE_FULL			(GODB_ERROR - 7)		/** Page ran out of space - internal error */
#define GODB_MAP_RESIZED		(GODB_ERROR - 8)		/** Database contents grew benyond environment mapsize */
#define GODB_INCOMPATIBLE		(GODB_ERROR - 9)		/** Database flags changes (or would change) */
/** @} */

/** Map for errno-related constants */
#define GODB_ERRNO_MAP(XX)											\
	/* No error */													\
	XX(GODB_SUCCESS, 		"success")								\
	XX(GODB_ERROR, 			"generic error")						\
																	
	/* DB Create Errors */											
	//XX(CORRUPTED, "db file corrupted")								
	/* Other error types */
	// XX(INVALID_FILE, "invalid database file")	

/* Define HPE_* values for each errno value above */
#define GODB_ERRNO_GEN(n, s) HPE_##n,
enum godb_errno {
	GODB_ERRNO_MAP(GODB_ERRNO_GEN)
};
#undef GODB_ERRNO_GEN

/* Get a database errno value */
#define GODB_DATABASE_ERRNO(p)			((enum godb_errno) (p)->godb_errno)




	
	
/* Move these to operation level flags */
	
#define GODB_RESERVE 0x0001/* Reserves a space for a data chunk and returns a pointer to it */
#define GODB_MULTIPLE 0x0001/* Writing batch items in one call */


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
#define GET_PAGESIZE(x) 							((x) = sysconf(_SC_PAGE_SIZE))
#define GODB_FLUSH_DATA_SYNC(file_descriptor) 		(!FlushFileBuffers(fd))
#define GODB_MSYNC(addr, len, flags) 				(!FlushViewOfFile(addr, len))
#define GODB_CLOSE_FILE(file_descriptor)			(CloseHandle(file_descriptor) ? 0 : -1)
#define GODB_MEM_UNMAP(ptr, len)					UnmapViewOfFile(ptr)

#ifdef O_CLOEXEC /* Linux: Open file and set FD_CLOEXEC atomically */
#	define GODB_CLOEXEC		O_CLOEXEC
#else
	 int fdflags;
#	define GODB_CLOEXEC		0
#endif

#define READ_WRITE_CREATE		O_RDWR | O_CREAT | GODB_CLOEXEC




typedef struct GODB_database GODB_database;

HANDLE godb_open_file(char *path, int file_mode);

/** Flag values for godb_database.flags field */
enum godb_flags
	{	F_OPEN_EXISTING			= 1 << 0
	,	F_CREATE_IF_MISSING		= 1 << 1
	,	F_READ_ONLY				= 1 << 2	/* Open database for read only access */
	,	F_NO_SYNC				= 1 << 3 /* Don't fsync after commit */
	,	F_NO_SYNC_HEADER		= 1 << 4 /* Don't fsync when writing to header metadata */
	};




//typedef struct godb_file godb_file;
//struct godb_file {
//	const struct godb_io_methods *pMethods; /* Methods for an open db file */
//};

//typedef struct godb_io_methods godb_io_methods;
//struct godb_io_methods {

//};
/** GODB Database root structure */
//struct GODB_database {
	/** PRIVATE **/

	/** READ-ONLY **/

	/** PUBLIC **/

//};


//int godb_open(GODB_database *db);
//int godb_get_flags(GODB_database *db, godb_flags *ret_flags);
//int godb_get_path(GODB_database *db, const char **ret_path);



#endif

/** @brief Create an MDB environment handle.
	 *
	 * This function allocates memory for a #MDB_env structure. To release
	 * the allocated memory and discard the handle, call #mdb_env_close().
	 * Before the handle may be used, it must be opened using #mdb_env_open().
	 * Various other options may also need to be set before opening the handle,
	 * e.g. #mdb_env_set_mapsize(), #mdb_env_set_maxreaders(), #mdb_env_set_maxdbs(),
	 * depending on usage requirements.
	 * @param[out] env The address where the new handle will be stored
	 * @return A non-zero error value on failure and 0 on success.
	 */
//int  godb_env_create(GODB_database **db);

	/** @brief Open an environment handle.
	 *
	 * If this function fails, #mdb_env_close() must be called to discard the #MDB_env handle.
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[in] path The directory in which the database files reside. This
	 * directory must already exist and be writable.
	 * @param[in] flags Special options for this environment. This parameter
	 * must be set to 0 or by bitwise OR'ing together one or more of the
	 * values described here.
	 * Flags set by mdb_env_set_flags() are also used.
	 * <ul>
	 *	<li>#MDB_FIXEDMAP
	 *      use a fixed address for the mmap region. This flag must be specified
	 *      when creating the environment, and is stored persistently in the environment.
	 *		If successful, the memory map will always reside at the same virtual address
	 *		and pointers used to reference data items in the database will be constant
	 *		across multiple invocations. This option may not always work, depending on
	 *		how the operating system has allocated memory to shared libraries and other uses.
	 *		The feature is highly experimental.
	 *	<li>#MDB_NOSUBDIR
	 *		By default, MDB creates its environment in a directory whose
	 *		pathname is given in \b path, and creates its data and lock files
	 *		under that directory. With this option, \b path is used as-is for
	 *		the database main data file. The database lock file is the \b path
	 *		with "-lock" appended.
	 *	<li>#MDB_RDONLY
	 *		Open the environment in read-only mode. No write operations will be
	 *		allowed. MDB will still modify the lock file - except on read-only
	 *		filesystems, where MDB does not use locks.
	 *	<li>#MDB_WRITEMAP
	 *		Use a writeable memory map unless MDB_RDONLY is set. This is faster
	 *		and uses fewer mallocs, but loses protection from application bugs
	 *		like wild pointer writes and other bad updates into the database.
	 *		Incompatible with nested transactions.
	 *	<li>#MDB_NOMETASYNC
	 *		Flush system buffers to disk only once per transaction, omit the
	 *		metadata flush. Defer that until the system flushes files to disk,
	 *		or next non-MDB_RDONLY commit or #mdb_env_sync(). This optimization
	 *		maintains database integrity, but a system crash may undo the last
	 *		committed transaction. I.e. it preserves the ACI (atomicity,
	 *		consistency, isolation) but not D (durability) database property.
	 *		This flag may be changed at any time using #mdb_env_set_flags().
	 *	<li>#MDB_NOSYNC
	 *		Don't flush system buffers to disk when committing a transaction.
	 *		This optimization means a system crash can corrupt the database or
	 *		lose the last transactions if buffers are not yet flushed to disk.
	 *		The risk is governed by how often the system flushes dirty buffers
	 *		to disk and how often #mdb_env_sync() is called.  However, if the
	 *		filesystem preserves write order and the #MDB_WRITEMAP flag is not
	 *		used, transactions exhibit ACI (atomicity, consistency, isolation)
	 *		properties and only lose D (durability).  I.e. database integrity
	 *		is maintained, but a system crash may undo the final transactions.
	 *		Note that (#MDB_NOSYNC | #MDB_WRITEMAP) leaves the system with no
	 *		hint for when to write transactions to disk, unless #mdb_env_sync()
	 *		is called. (#MDB_MAPASYNC | #MDB_WRITEMAP) may be preferable.
	 *		This flag may be changed at any time using #mdb_env_set_flags().
	 *	<li>#MDB_MAPASYNC
	 *		When using #MDB_WRITEMAP, use asynchronous flushes to disk.
	 *		As with #MDB_NOSYNC, a system crash can then corrupt the
	 *		database or lose the last transactions. Calling #mdb_env_sync()
	 *		ensures on-disk database integrity until next commit.
	 *		This flag may be changed at any time using #mdb_env_set_flags().
	 *	<li>#MDB_NOTLS
	 *		Don't use Thread-Local Storage. Tie reader locktable slots to
	 *		#MDB_txn objects instead of to threads. I.e. #mdb_txn_reset() keeps
	 *		the slot reseved for the #MDB_txn object. A thread may use parallel
	 *		read-only transactions. A read-only transaction may span threads if
	 *		the user synchronizes its use. Applications that multiplex many
	 *		user threads over individual OS threads need this option. Such an
	 *		application must also serialize the write transactions in an OS
	 *		thread, since MDB's write locking is unaware of the user threads.
	 * </ul>
	 * @param[in] mode The UNIX permissions to set on created files. This parameter
	 * is ignored on Windows.
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>#MDB_VERSION_MISMATCH - the version of the MDB library doesn't match the
	 *	version that created the database environment.
	 *	<li>#MDB_INVALID - the environment file headers are corrupted.
	 *	<li>ENOENT - the directory specified by the path parameter doesn't exist.
	 *	<li>EACCES - the user didn't have permission to access the environment files.
	 *	<li>EAGAIN - the environment was locked by another process.
	 * </ul>
	 */
//int  godb_env_open(GODB_database *db, const char *path, unsigned int flags, godb_mode_t mode);

	/** @brief Copy an MDB environment to the specified path.
	 *
	 * This function may be used to make a backup of an existing environment.
	 * @param[in] env An environment handle returned by #mdb_env_create(). It
	 * must have already been opened successfully.
	 * @param[in] path The directory in which the copy will reside. This
	 * directory must already exist and be writable but must otherwise be
	 * empty.
	 * @return A non-zero error value on failure and 0 on success.
	 */
//int  godb_env_copy(GODB_database *db, const char *path);

	/** @brief Copy an MDB environment to the specified file descriptor.
	 *
	 * This function may be used to make a backup of an existing environment.
	 * @param[in] env An environment handle returned by #mdb_env_create(). It
	 * must have already been opened successfully.
	 * @param[in] fd The filedescriptor to write the copy to. It must
	 * have already been opened for Write access.
	 * @return A non-zero error value on failure and 0 on success.
	 */
//int  godb_env_copyfd(GODB_database *db, godb_filehandle_t fd);

	/** @brief Return statistics about the MDB environment.
	 *
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[out] stat The address of an #MDB_stat structure
	 * 	where the statistics will be copied
	 */
//int  godb_env_stat(GODB_database *db, GODB_stat *stat);

	/** @brief Return information about the MDB environment.
	 *
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[out] stat The address of an #MDB_envinfo structure
	 * 	where the information will be copied
	 */
//int  godb_env_info(GODB_database *db, GODB_envinfo *stat);

	/** @brief Flush the data buffers to disk.
	 *
	 * Data is always written to disk when #mdb_txn_commit() is called,
	 * but the operating system may keep it buffered. MDB always flushes
	 * the OS buffers upon commit as well, unless the environment was
	 * opened with #MDB_NOSYNC or in part #MDB_NOMETASYNC.
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[in] force If non-zero, force a synchronous flush.  Otherwise
	 *  if the environment has the #MDB_NOSYNC flag set the flushes
	 *	will be omitted, and with #MDB_MAPASYNC they will be asynchronous.
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified.
	 *	<li>EIO - an error occurred during synchronization.
	 * </ul>
	 */
//int  gdb_env_sync(GODB_database *db, int force);

	/** @brief Close the environment and release the memory map.
	 *
	 * Only a single thread may call this function. All transactions, databases,
	 * and cursors must already be closed before calling this function. Attempts to
	 * use any such handles after calling this function will cause a SIGSEGV.
	 * The environment handle will be freed and must not be used again after this call.
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 */
//void godb_env_close(GODB_database *db);

	/** @brief Set environment flags.
	 *
	 * This may be used to set some flags in addition to those from
	 * #mdb_env_open(), or to unset these flags.
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[in] flags The flags to change, bitwise OR'ed together
	 * @param[in] onoff A non-zero value sets the flags, zero clears them.
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//int  godb_env_set_flags(GODB_database *db, unsigned int flags, int onoff);
//int godb_set_flags(GODB_database *db, unsigned int flags, unsigned short turnOn);

	/** @brief Get environment flags.
	 *
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[out] flags The address of an integer to store the flags
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//int  godb_get_flags(GODB_database *db, unsigned int *ret_flags);

	/** @brief Return the path that was used in #mdb_env_open().
	 *
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[out] path Address of a string pointer to contain the path. This
	 * is the actual string in the environment, not a copy. It should not be
	 * altered in any way.
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//int  godb_get_path(GODB_database *db, const char **ret_path);

	/** @brief Set the size of the memory map to use for this environment.
	 *
	 * The size should be a multiple of the OS page size. The default is
	 * 10485760 bytes. The size of the memory map is also the maximum size
	 * of the database. The value should be chosen as large as possible,
	 * to accommodate future growth of the database.
	 * This function may only be called after #mdb_env_create() and before #mdb_env_open().
	 * The size may be changed by closing and reopening the environment.
	 * Any attempt to set a size smaller than the space already consumed
	 * by the environment will be silently changed to the current size of the used space.
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[in] size The size in bytes
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified, or the environment is already open.
	 * </ul>
	 */
//int  godb_env_set_mapsize(GODB_database *db, size_t size);
