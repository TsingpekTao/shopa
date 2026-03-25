package packed

import (
	// MySQL driver is required for gdb to create a mysql connection from the link string.
	// We keep it as a blank import so it registers itself without polluting packages.
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)
