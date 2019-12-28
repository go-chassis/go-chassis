// env only contains env, no other params
package env

import "os"

var (
	ChassisHome = os.Getenv("CHASSIS_HOME")
)
