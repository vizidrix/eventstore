#ifndef ES_H_INCLUDED
#define ES_H_INCLUDED

#ifdef __cplusplus
extern "C"
#endif

#include <errno.h>
#include <stddef.h>

/*  Handle DSO symbol visibility                                             */
#if defined _WIN32
#   if defined ES_EXPORTS
#       define ES_EXPORT __declspec(dllexport)
#   else
#       define ES_EXPORT __declspec(dllimport)
#   endif
#else
#   if defined __SUNPRO_C
#       define ES_EXPORT __global
#   elif (defined __GNUC__ && __GNUC__ >= 4) || \
          defined __INTEL_COMPILER || defined __clang__
#       define ES_EXPORT __attribute__ ((visibility("default")))
#   else
#       define ES_EXPORT
#   endif
#endif

/******************************************************************************/
/*  ABI versioning support.                                                   */
/******************************************************************************/

/*  Don't change this unless you know exactly what you're doing and have      */
/*  read and understand the following documents:                              */
/*  www.gnu.org/software/libtool/manual/html_node/Libtool-versioning.html     */
/*  www.gnu.org/software/libtool/manual/html_node/Updating-version-info.html  */

/*  The current interface version. */
#define ES_VERSION_CURRENT 0

/*  The latest revision of the current interface. */
#define ES_VERSION_REVISION 0

/*  How many past interface versions are still supported. */
#define ES_VERSION_AGE 0

/******************************************************************************/
/*  Zero-copy support.                                                        */
/******************************************************************************/

#define ES_MSG ((size_t) -1)

ES_EXPORT void *es_alloc_batch (
	uint32_t domain_id, 		/* Application partition */
	uint32_t kind_id, 			/* Kind partition (aggregate type) */
	uint64_t aggregate_id,		/* Aggregate instance ids (event partition) */
	char count);				/* Number of contiguous buffers to allocate */
ES_EXPORT int es_free_batch (void *batch);

/******************************************************************************/
/*  Type definition.                                                          */
/******************************************************************************/

#define WRITER_HANDLE	int

struct es_batch {};

ES_EXPORT WRITER_HANDLE es_open_write(const char *path);
ES_EXPORT int es_close_write(WRITER_HANDLE db);
ES_EXPORT int es_publish (WRITER_HANDLE db, const struct es_batch *batch, int flags);
ES_EXPORT int es_cancel (WRITER_HANDLE db, const struct es_batch *batch);

#undef ES_EXPORT

#ifdef __cplusplus
}
#endif

#endif