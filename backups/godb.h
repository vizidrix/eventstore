

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
//#include <errno.h>
//#include <limits.h>
//#include <stddef.h>
//#include <inttypes.h>
//#include <stdio.h>

//#include <string.h>
////#include <time.h>
//#include <unistd.h>


/*
typedef struct {
	char*	data_path;
} GODB_config;

typedef struct {
	unsigned short		event_type;
} GODB_header;

typedef struct {
	char*	body; // Should pointer to mem mapped be void*?
	int 	length;
} GODB_value;
*/
	

//typedef	mode_t	godb_mode_t;

	/** An abstraction for a file handle. */
//typedef int godb_filehandle_t;









/** @brief A handle for an individual database in the DB environment. */
//typedef unsigned int	GODB_dbi;

/** @brief Opaque structure for navigating through a database */
//typedef struct GODB_cursor GODB_cursor;

/** @brief Generic structure used for passing keys and data in and out
 * of the database.
 *
 * Values returned from the database are valid indefinately as event store is immutable
 */
//typedef struct GODB_value {
	//size_t		 mv_size;	/**< size of the data item */
//	unsigned int		mv_size;		/**< size of the data item */
//	void				*mv_data;		/**< address of the data item */
//} GODB_value;

/** @breif Structure to capture runtime statistics */
//typedef struct GODB_stat {
//	unsigned int ms_psize; 		/**< Size of a database page. */

//} GODB_stat;

/** @breif Information about the environment */
//typedef struct GODB_envinfo {
//	size_t	me_mapsize; 		/**< Size of the data memory map. */
//} GODB_envinfo;


//int *create_file(char *path, int file_mode, HANDLE file_handle);

	/** @brief Return the mdb library version information.
	 *
	 * @param[out] major if non-NULL, the library major version number is copied here
	 * @param[out] minor if non-NULL, the library minor version number is copied here
	 * @param[out] patch if non-NULL, the library patch version number is copied here
	 * @retval "version string" The library version as a string
	 */
//char *godb_version(int *major, int *minor, int *patch);

	/** @brief Return a string describing a given error code.
	 *
	 * This function is a superset of the ANSI C X3.159-1989 (ANSI C) strerror(3)
	 * function. If the error code is greater than or equal to 0, then the string
	 * returned by the system function strerror(3) is returned. If the error code
	 * is less than 0, an error string corresponding to the MDB library error is
	 * returned. See @ref errors for a list of MDB-specific error codes.
	 * @param[in] err The error code
	 * @retval "error message" The description of the error
	 */
//char *godb_strerror(int err);


	/** @brief Open a database in the environment.
	 *
	 * A database handle denotes the name and parameters of a database,
	 * independently of whether such a database exists.
	 * The database handle may be discarded by calling #mdb_dbi_close().
	 * The old database handle is returned if the database was already open.
	 * The handle must only be closed once.
	 * The database handle will be private to the current transaction until
	 * the transaction is successfully committed. If the transaction is
	 * aborted the handle will be closed automatically.
	 * After a successful commit the
	 * handle will reside in the shared environment, and may be used
	 * by other transactions. This function must not be called from
	 * multiple concurrent transactions. A transaction that uses this function
	 * must finish (either commit or abort) before any other transaction may
	 * use this function.
	 *
	 * To use named databases (with name != NULL), #mdb_env_set_maxdbs()
	 * must be called before opening the environment.
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] name The name of the database to open. If only a single
	 * 	database is needed in the environment, this value may be NULL.
	 * @param[in] flags Special options for this database. This parameter
	 * must be set to 0 or by bitwise OR'ing together one or more of the
	 * values described here.
	 * <ul>
	 *	<li>#MDB_REVERSEKEY
	 *		Keys are strings to be compared in reverse order, from the end
	 *		of the strings to the beginning. By default, Keys are treated as strings and
	 *		compared from beginning to end.
	 *	<li>#MDB_DUPSORT
	 *		Duplicate keys may be used in the database. (Or, from another perspective,
	 *		keys may have multiple data items, stored in sorted order.) By default
	 *		keys must be unique and may have only a single data item.
	 *	<li>#MDB_INTEGERKEY
	 *		Keys are binary integers in native byte order. Setting this option
	 *		requires all keys to be the same size, typically sizeof(int)
	 *		or sizeof(size_t).
	 *	<li>#MDB_DUPFIXED
	 *		This flag may only be used in combination with #MDB_DUPSORT. This option
	 *		tells the library that the data items for this database are all the same
	 *		size, which allows further optimizations in storage and retrieval. When
	 *		all data items are the same size, the #MDB_GET_MULTIPLE and #MDB_NEXT_MULTIPLE
	 *		cursor operations may be used to retrieve multiple items at once.
	 *	<li>#MDB_INTEGERDUP
	 *		This option specifies that duplicate data items are also integers, and
	 *		should be sorted as such.
	 *	<li>#MDB_REVERSEDUP
	 *		This option specifies that duplicate data items should be compared as
	 *		strings in reverse order.
	 *	<li>#MDB_CREATE
	 *		Create the named database if it doesn't exist. This option is not
	 *		allowed in a read-only transaction or a read-only environment.
	 * </ul>
	 * @param[out] dbi Address where the new #MDB_dbi handle will be stored
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>#MDB_NOTFOUND - the specified database doesn't exist in the environment
	 *		and #MDB_CREATE was not specified.
	 *	<li>#MDB_DBS_FULL - too many databases have been opened. See #mdb_env_set_maxdbs().
	 * </ul>
	 */
//----int  mdb_dbi_open(MDB_txn *txn, const char *name, unsigned int flags, MDB_dbi *dbi);

	/** @brief Retrieve statistics for a database.
	 *
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 * @param[out] stat The address of an #MDB_stat structure
	 * 	where the statistics will be copied
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_stat(MDB_txn *txn, MDB_dbi dbi, MDB_stat *stat);

	/** @brief Retrieve the DB flags for a database handle.
	 *
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 * @param[out] flags Address where the flags will be returned.
	 * @return A non-zero error value on failure and 0 on success.
	 */
//----int mdb_dbi_flags(MDB_env *env, MDB_dbi dbi, unsigned int *flags);

	/** @brief Close a database handle.
	 *
	 * This call is not mutex protected. Handles should only be closed by
	 * a single thread, and only if no other threads are going to reference
	 * the database handle or one of its cursors any further. Do not close
	 * a handle if an existing transaction has modified its database.
	 * @param[in] env An environment handle returned by #mdb_env_create()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 */
//----void mdb_dbi_close(MDB_env *env, MDB_dbi dbi);

	/** @brief Delete a database and/or free all its pages.
	 *
	 * If the \b del parameter is 1, the DB handle will be closed
	 * and the DB will be deleted.
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 * @param[in] del 1 to delete the DB from the environment,
	 * 0 to just free its pages.
	 * @return A non-zero error value on failure and 0 on success.
	 */
//----int  mdb_drop(MDB_txn *txn, MDB_dbi dbi, int del);


	/** @brief Get items from a database.
	 *
	 * This function retrieves key/data pairs from the database. The address
	 * and length of the data associated with the specified \b key are returned
	 * in the structure to which \b data refers.
	 * If the database supports duplicate keys (#MDB_DUPSORT) then the
	 * first data item for the key will be returned. Retrieval of other
	 * items requires the use of #mdb_cursor_get().
	 *
	 * @note The memory pointed to by the returned values is owned by the
	 * database. The caller need not dispose of the memory, and may not
	 * modify it in any way. For values returned in a read-only transaction
	 * any modification attempts will cause a SIGSEGV.
	 * @note Values returned from the database are valid only until a
	 * subsequent update operation, or the end of the transaction.
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 * @param[in] key The key to search for in the database
	 * @param[out] data The data corresponding to the key
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>#MDB_NOTFOUND - the key was not in the database.
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_get(MDB_txn *txn, MDB_dbi dbi, MDB_val *key, MDB_val *data);

	/** @brief Store items into a database.
	 *
	 * This function stores key/data pairs in the database. The default behavior
	 * is to enter the new key/data pair, replacing any previously existing key
	 * if duplicates are disallowed, or adding a duplicate data item if
	 * duplicates are allowed (#MDB_DUPSORT).
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 * @param[in] key The key to store in the database
	 * @param[in,out] data The data to store
	 * @param[in] flags Special options for this operation. This parameter
	 * must be set to 0 or by bitwise OR'ing together one or more of the
	 * values described here.
	 * <ul>
	 *	<li>#MDB_NODUPDATA - enter the new key/data pair only if it does not
	 *		already appear in the database. This flag may only be specified
	 *		if the database was opened with #MDB_DUPSORT. The function will
	 *		return #MDB_KEYEXIST if the key/data pair already appears in the
	 *		database.
	 *	<li>#MDB_NOOVERWRITE - enter the new key/data pair only if the key
	 *		does not already appear in the database. The function will return
	 *		#MDB_KEYEXIST if the key already appears in the database, even if
	 *		the database supports duplicates (#MDB_DUPSORT). The \b data
	 *		parameter will be set to point to the existing item.
	 *	<li>#MDB_RESERVE - reserve space for data of the given size, but
	 *		don't copy the given data. Instead, return a pointer to the
	 *		reserved space, which the caller can fill in later - before
	 *		the next update operation or the transaction ends. This saves
	 *		an extra memcpy if the data is being generated later.
	 *	<li>#MDB_APPEND - append the given key/data pair to the end of the
	 *		database. No key comparisons are performed. This option allows
	 *		fast bulk loading when keys are already known to be in the
	 *		correct order. Loading unsorted keys with this flag will cause
	 *		data corruption.
	 *	<li>#MDB_APPENDDUP - as above, but for sorted dup data.
	 * </ul>
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>#MDB_MAP_FULL - the database is full, see #mdb_env_set_mapsize().
	 *	<li>#MDB_TXN_FULL - the transaction has too many dirty pages.
	 *	<li>EACCES - an attempt was made to write in a read-only transaction.
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_put(MDB_txn *txn, MDB_dbi dbi, MDB_val *key, MDB_val *data,
//----			    unsigned int flags);

	/** @brief Delete items from a database.
	 *
	 * This function removes key/data pairs from the database.
	 * If the database does not support sorted duplicate data items
	 * (#MDB_DUPSORT) the data parameter is ignored.
	 * If the database supports sorted duplicates and the data parameter
	 * is NULL, all of the duplicate data items for the key will be
	 * deleted. Otherwise, if the data parameter is non-NULL
	 * only the matching data item will be deleted.
	 * This function will return #MDB_NOTFOUND if the specified key/data
	 * pair is not in the database.
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 * @param[in] key The key to delete from the database
	 * @param[in] data The data to delete
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EACCES - an attempt was made to write in a read-only transaction.
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_del(MDB_txn *txn, MDB_dbi dbi, MDB_val *key, MDB_val *data);

	/** @brief Create a cursor handle.
	 *
	 * A cursor is associated with a specific transaction and database.
	 * A cursor cannot be used when its database handle is closed.  Nor
	 * when its transaction has ended, except with #mdb_cursor_renew().
	 * It can be discarded with #mdb_cursor_close().
	 * A cursor in a write-transaction can be closed before its transaction
	 * ends, and will otherwise be closed when its transaction ends.
	 * A cursor in a read-only transaction must be closed explicitly, before
	 * or after its transaction ends. It can be reused with
	 * #mdb_cursor_renew() before finally closing it.
	 * @note Earlier documentation said that cursors in every transaction
	 * were closed when the transaction committed or aborted.
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] dbi A database handle returned by #mdb_dbi_open()
	 * @param[out] cursor Address where the new #MDB_cursor handle will be stored
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_cursor_open(MDB_txn *txn, MDB_dbi dbi, MDB_cursor **cursor);

	/** @brief Close a cursor handle.
	 *
	 * The cursor handle will be freed and must not be used again after this call.
	 * Its transaction must still be live if it is a write-transaction.
	 * @param[in] cursor A cursor handle returned by #mdb_cursor_open()
	 */
//----void mdb_cursor_close(MDB_cursor *cursor);

	/** @brief Renew a cursor handle.
	 *
	 * A cursor is associated with a specific transaction and database.
	 * Cursors that are only used in read-only
	 * transactions may be re-used, to avoid unnecessary malloc/free overhead.
	 * The cursor may be associated with a new read-only transaction, and
	 * referencing the same database handle as it was created with.
	 * This may be done whether the previous transaction is live or dead.
	 * @param[in] txn A transaction handle returned by #mdb_txn_begin()
	 * @param[in] cursor A cursor handle returned by #mdb_cursor_open()
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_cursor_renew(MDB_txn *txn, MDB_cursor *cursor);

/** @brief Return the cursor's database handle.
	 *
	 * @param[in] cursor A cursor handle returned by #mdb_cursor_open()
	 */
//----MDB_dbi mdb_cursor_dbi(MDB_cursor *cursor);


	/** @brief Retrieve by cursor.
	 *
	 * This function retrieves key/data pairs from the database. The address and length
	 * of the key are returned in the object to which \b key refers (except for the
	 * case of the #MDB_SET option, in which the \b key object is unchanged), and
	 * the address and length of the data are returned in the object to which \b data
	 * refers.
	 * See #mdb_get() for restrictions on using the output values.
	 * @param[in] cursor A cursor handle returned by #mdb_cursor_open()
	 * @param[in,out] key The key for a retrieved item
	 * @param[in,out] data The data of a retrieved item
	 * @param[in] op A cursor operation #MDB_cursor_op
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>#MDB_NOTFOUND - no matching key found.
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_cursor_get(MDB_cursor *cursor, MDB_val *key, MDB_val *data,
//----			    MDB_cursor_op op);

	/** @brief Store by cursor.
	 *
	 * This function stores key/data pairs into the database.
	 * If the function fails for any reason, the state of the cursor will be
	 * unchanged. If the function succeeds and an item is inserted into the
	 * database, the cursor is always positioned to refer to the newly inserted item.
	 * @param[in] cursor A cursor handle returned by #mdb_cursor_open()
	 * @param[in] key The key operated on.
	 * @param[in] data The data operated on.
	 * @param[in] flags Options for this operation. This parameter
	 * must be set to 0 or one of the values described here.
	 * <ul>
	 *	<li>#MDB_CURRENT - overwrite the data of the key/data pair to which
	 *		the cursor refers with the specified data item. The \b key
	 *		parameter is ignored.
	 *	<li>#MDB_NODUPDATA - enter the new key/data pair only if it does not
	 *		already appear in the database. This flag may only be specified
	 *		if the database was opened with #MDB_DUPSORT. The function will
	 *		return #MDB_KEYEXIST if the key/data pair already appears in the
	 *		database.
	 *	<li>#MDB_NOOVERWRITE - enter the new key/data pair only if the key
	 *		does not already appear in the database. The function will return
	 *		#MDB_KEYEXIST if the key already appears in the database, even if
	 *		the database supports duplicates (#MDB_DUPSORT).
	 *	<li>#MDB_RESERVE - reserve space for data of the given size, but
	 *		don't copy the given data. Instead, return a pointer to the
	 *		reserved space, which the caller can fill in later. This saves
	 *		an extra memcpy if the data is being generated later.
	 *	<li>#MDB_APPEND - append the given key/data pair to the end of the
	 *		database. No key comparisons are performed. This option allows
	 *		fast bulk loading when keys are already known to be in the
	 *		correct order. Loading unsorted keys with this flag will cause
	 *		data corruption.
	 *	<li>#MDB_APPENDDUP - as above, but for sorted dup data.
	 *	<li>#MDB_MULTIPLE - store multiple contiguous data elements in a
	 *		single request. This flag may only be specified if the database
	 *		was opened with #MDB_DUPFIXED. The \b data argument must be an
	 *		array of two MDB_vals. The mv_size of the first MDB_val must be
	 *		the size of a single data element. The mv_data of the first MDB_val
	 *		must point to the beginning of the array of contiguous data elements.
	 *		The mv_size of the second MDB_val must be the count of the number
	 *		of data elements to store. On return this field will be set to
	 *		the count of the number of elements actually written. The mv_data
	 *		of the second MDB_val is unused.
	 * </ul>
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>#MDB_MAP_FULL - the database is full, see #mdb_env_set_mapsize().
	 *	<li>#MDB_TXN_FULL - the transaction has too many dirty pages.
	 *	<li>EACCES - an attempt was made to modify a read-only database.
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_cursor_put(MDB_cursor *cursor, MDB_val *key, MDB_val *data,
//----				unsigned int flags);

	/** @brief Delete current key/data pair
	 *
	 * This function deletes the key/data pair to which the cursor refers.
	 * @param[in] cursor A cursor handle returned by #mdb_cursor_open()
	 * @param[in] flags Options for this operation. This parameter
	 * must be set to 0 or one of the values described here.
	 * <ul>
	 *	<li>#MDB_NODUPDATA - delete all of the data items for the current key.
	 *		This flag may only be specified if the database was opened with #MDB_DUPSORT.
	 * </ul>
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EACCES - an attempt was made to modify a read-only database.
	 *	<li>EINVAL - an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_cursor_del(MDB_cursor *cursor, unsigned int flags);

	/** @brief Return count of duplicates for current key.
	 *
	 * This call is only valid on databases that support sorted duplicate
	 * data items #MDB_DUPSORT.
	 * @param[in] cursor A cursor handle returned by #mdb_cursor_open()
	 * @param[out] countp Address where the count will be stored
	 * @return A non-zero error value on failure and 0 on success. Some possible
	 * errors are:
	 * <ul>
	 *	<li>EINVAL - cursor is not initialized, or an invalid parameter was specified.
	 * </ul>
	 */
//----int  mdb_cursor_count(MDB_cursor *cursor, size_t *countp);


	/** @brief A callback function used to print a message from the library.
	 *
	 * @param[in] msg The string to be printed.
	 * @param[in] ctx An arbitrary context pointer for the callback.
	 * @return < 0 on failure, 0 on success.
	 */
//----typedef int (MDB_msg_func)(const char *msg, void *ctx);

