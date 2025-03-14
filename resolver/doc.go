/*
package resolver provides an implementation for resolving parameterized strings.

The resolver is used to convert tokens into actual values.  The resolver will
handle any token-recursion it encounters; so for example it can resolve a property
token that itself has an environment token:

---

	ResolverConfig.Properties: {"foo": "${env:bar}"}
	environment: {"bar" = "hello"}
	Resolve("${prop:foo}") -> "hello"

---

Tokens are specified using '${type:value}' notation.

There are several supported token types, and the value specifcation depends on its type.
If Resolver cannot resolve the value of a token, then the token is left unchanged.

If "type" is omitted, then the resolver assumes a property token.
Therefore, "${prop: foo}" and "${foo}" are identical.

Resolver supports the following token types:

	"prop"  : "value" is the name of a property found in the Config.Properties map.
		  During resolution, the token is replaced with the value of the named property.
	"env"   : "value" is the name of an environment variable.
		  During resolution, the token is replaced with the value of the environment variable.
	"date"|"time"|"datetime" : "value" specifies a date, time, or date+time expression.
	"epoch" : "value" is the number of seconds since the Unix epoch (January 1, 1970 UTC)

	For date/time resolution, "value" is expressed using the format:
		"date[.(format)[ <+|-> delta<S>]]"
	Where:
		"date" is either the literal "now", or the date/time according to the "format" specification
		"format" is the date/time format specification; this is used to parse literal dates and also to stringfy the result.
			If not otherwise specified, the "RFC3339" format is used.  See "Format Specifiers" below for more information.
		"[<+|-> delta<S>]" provides a means of adding or subtracting time from the "date":
			"<+|->" indicates addition or subraction.
				This is REQUIRED when expressing a time/date delta and must be either "+" or "-" (but not both).
			"delta" is an integer specifying the magnitude of the change
			"<S>" indicates the scope of the change, which must be one of the following case-sensitive values [YMWDhms],
			which correspond to Year, Month, Week, Day, hour, minute, and second respectively.

	For "epoch", the format specification is ignored.  You may specify a delta for hours, minutes, or seconds.

For Example:

Assuming the current date and time ("now") is Thursday July 6, 2025 3:45:21 pm, EST:

	"${date:now}" -> "2025-07-06"
	"${time:now}" -> "15:45:21Z"
	"${datetime:now}" -> "2025-07-06T15:45:21Z"
	"${date:now(RFC850)}" -> "Thursday, 06-Jul-25 15:45:21 EST"
	"${date:now + 1Y}" -> "2026-07-10"  (now plus one year)
	"${date:now + 2M}" -> "2025-09-06"  (now plus 2 months)
	"${date:now - 4D}" -> "2025-07-02"  (now less 4 days)
	"${time:now + 3h}" -> "18:45:21"    (now plus 3 hours)
	"${epoch:now + 1h}" -> "1751820321" (now plus 1 hour)
	... and so on.

Format Specifiers:

Date/Time formats can be specified in three ways:

1. Using the name of one of the Layout constants from the time package ("RFC3339", "DateOnly", etc.)

2. Using the Layout value specification documented in the time package (refer to the time.Layout constant)

3. Using "YYYYMMDDhhmmss" case-sensitive date formatting strings, where:

	"Year", "YYYY", "YY"         : resolve to "2025", "2025", and "25" respectively
	"Month", "Mon", "MM", "_M"   : resolve to "July", "Jul", "07", and "7" respectively
	"Weekday", "Day", "DD", "_D" : resolve to "Thursday", "Thu", "06" and "6" respectively
	"HH", "hh" "_h"              : resolve to "15", "03" and "3" respetively (lowercase values indicate 12h notation)
	"mm", "ss"                   : resolve to "45", "21"
	"AM", "PM"                   : resolve to either "AM" or "PM", depending on the hour
	                               (typically only needed with 12h notation)
	"TimeZone", "TZ", "ZZ"       : resolve to "EST", "-0400", and "Z04:00" respectively
*/
package resolver
