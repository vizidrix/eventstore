/******************************************************************************
 *
 * Function Name:  manualCopy
 *
 * Description: Manually copies data from memory to memory.  This is used by
 * sysFastMemCopy to copy a few lingering bytes at the beginning and end.
 *
 *****************************************************************************/
 
inline void manualCopy( uint8 *pDest, uint8 *pSrc, uint32 len )
{
    uint32 i;
 
    // Manually copy the data
    for ( i = 0; i < len; i++ )
    {
        // Copy data from source to destination
        *pDest++ = *pSrc++;
    }
}
 
/******************************************************************************
 *
 * Function Name:  sysFastMemCopy
 *
 * Description: Assuming that your processor can do 32-bit memory accesses
 * and contains a barrel shifter and that you are using an efficient
 * compiler, then this memory-to-memory copy procedure will probably be more
 * efficient than just using the traditional memcpy procedure if the number
 * of bytes to copy is greater than about 20.  It works by doing 32-bit
 * reads and writes instead of using 8-bit memory accesses.
 *
 * NOTE that this procedure assumes a Little Endian processor!  The shift
 * operators ">>" and "<<" should all be reversed for Big Endian.
 *
 * NEVER use this when the number of bytes to be copied is less than about
 * 10, since it may not work for a small number of bytes.  Also, do not use
 * this when the source and destination regions overlap.
 *
 * NOTE that this may NOT be faster than memcpy if your processor supports a
 * really fast cache memory!
 *
 * Timing for this sysFastMemCopy varies some according to which shifts need
 * to be done.  The following results are from one attempt to measure timing
 * on a Cortex M4 processor running at 48 MHz.
 *
 *                           MEMCPY        FAST
 *                  BYTES  bytes/usec   bytes/usec
 *                  -----  ----------  ------------
 *                    50       4.2      6.3 to  6.3
 *                   100       4.5      8.3 to 10.0
 *                   150       4.8     10.0 to 11.5
 *                   200       4.9     10.5 to 12.5
 *                   250       5.1     11.4 to 13.2
 *                   300       5.1     11.5 to 13.6
 *                   350       5.1     12.1 to 14.6
 *                   400       5.1     12.1 to 14.8
 *                   450       5.2     12.2 to 15.5
 *                   500       5.2     12.5 to 15.2
 *
 * The following macro can be used instead of memcpy to automatically select
 * whether to use memcpy or sysFastMemCopy:
 *
 *   #define MEMCOPY(pDst,pSrc,len) \
 *     (len) < 16 ? memcpy(pDst,pSrc,len) : sysFastMemCopy(pDst,pSrc,len);
 *
 *****************************************************************************/
 
void sysFastMemCopy( uint8 *pDest, uint8 *pSrc, uint32 len )
{
    uint32 srcCnt;
    uint32 destCnt;
    uint32 newLen;
    uint32 endLen;
    uint32 longLen;
    uint32 *pLongSrc;
    uint32 *pLongDest;
    uint32 longWord1;
    uint32 longWord2;
    uint32 methodSelect;
     
    // Determine the number of bytes in the first word of src and dest
    srcCnt = 4 - ( (uint32) pSrc & 0x03 );
    destCnt = 4 - ( (uint32) pDest & 0x03 );
     
    // Copy the initial bytes to the destination
    manualCopy( pDest, pSrc, destCnt );
     
    // Determine the number of bytes remaining
    newLen = len - destCnt;
     
    // Determine how many full long words to copy to the destination
    longLen = newLen / 4;
     
    // Determine number of lingering bytes to copy at the end
    endLen = newLen & 0x03;
     
    // Pick the initial long destination word to copy to
    pLongDest = (uint32*) ( pDest + destCnt );
     
    // Pick the initial source word to start our algorithm at
    if ( srcCnt <= destCnt )
    {
        // Advance to pSrc at the start of the next full word
        pLongSrc = (uint32*) ( pSrc + srcCnt );
    }
    else // There are still source bytes remaining in the first word
    {
        // Set pSrc to the start of the first full word
        pLongSrc = (uint32*) ( pSrc + srcCnt - 4 );
    }
     
    // There are 4 different longWord copy methods
    methodSelect = ( srcCnt - destCnt ) & 0x03;
     
    // Just copy one-to-one
    if ( methodSelect == 0 )
    {
        // Just copy the specified number of long words
        while ( longLen-- > 0 )
        {
            *pLongDest++ = *pLongSrc++;
        }
    }
    else if ( methodSelect == 1 )
    {
        // Get the first long word
        longWord1 = *pLongSrc++;
         
        // Copy words created by combining 2 adjacent long words
        while ( longLen-- > 0 )
        {
            // Get the next 32-bit word
            longWord2 = *pLongSrc++;
             
            // Write to the destination
            *pLongDest++ = ( longWord1 >> 24 ) | ( longWord2 << 8 );
             
            // Re-use the word just retrieved
            longWord1 = longWord2;
        }
    }
    else if ( methodSelect == 2 )
    {
        // Get the first long word
        longWord1 = *pLongSrc++;
         
        // Copy words created by combining 2 adjacent long words
        while ( longLen-- > 0 )
        {
            // Get the next 32-bit word
            longWord2 = *pLongSrc++;
             
            // Write to the destination
            *pLongDest++ = ( longWord1 >> 16 ) | ( longWord2 << 16 );
             
            // Re-use the word just retrieved
            longWord1 = longWord2;
        }
    }
    else // ( methodSelect == 3 )
    {
        // Get the first long word
        longWord1 = *pLongSrc++;
         
        // Copy words created by combining 2 adjacent long words
        while ( longLen-- > 0 )
        {
            // Get the next 32-bit word
            longWord2 = *pLongSrc++;
 
            // Write to the destination
            *pLongDest++ = ( longWord1 >> 8 ) | ( longWord2 << 24 );
 
            // Re-use the word just retrieved
            longWord1 = longWord2;
        }
    }
     
    // Copy any remaining bytes
    if ( endLen != 0 )
    {
        // The trailing bytes will be copied next
        pDest = (uint8*) pLongDest;
         
        // Determine where the trailing source bytes are located
        pSrc += len - endLen;
         
        // Copy the remaining bytes
        manualCopy( pDest, pSrc, endLen );
    }
}