package environment

import "time"

var ServerStartTime = time.Now()

func PassedTime() time.Duration {
	return time.Now().Sub(ServerStartTime)
}
