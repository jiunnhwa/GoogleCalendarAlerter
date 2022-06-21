/*

Assorted service functions.

*/
package service

// ObfuscateLast4Char replaces the last 4 char with *
func ObfuscateLast4Char(str string) string {
	return str[0:len(str)-4] + "****"
}
