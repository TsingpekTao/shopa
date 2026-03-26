package packed

import (
	// MySQL 驱动用于让 gdb 能根据连接串建立连接。
	// 以匿名导入方式注册驱动，避免污染其它包的命名空间。
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)
