/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package retry

import "time"

// Defaults
const (
	DefaultMaxAttempts = 5
	DefaultRetryDelay  = time.Second
)
