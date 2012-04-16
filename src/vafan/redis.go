// Copyright 2012 Saul Howard. All rights reserved.
//
// Redis.
//

package vafan

import (
	"github.com/fzzbt/radix"
)

var redisConfiguration = getRedisConf()

func getRedisConf() radix.Configuration {
    return radix.Configuration{
        Database: 0,  // (default: 0)
        Timeout:  10, // (default: 10)
        Address:  vafanConf.redis.address,
    }
}
