// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from 'react'
import classnames from 'classnames'

import Status from '@ttn-lw/components/status'

import DateTime from '@ttn-lw/lib/components/date-time'
import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import style from './last-seen.styl'

const computeDeltaInSeconds = (from, to) => {
  // Avoid situations when server clock is ahead of the browser clock.
  if (from > to) {
    return 0
  }

  return Math.floor((from - to) / 1000)
}

const LastSeen = React.forwardRef((props, ref) => {
  const { className, lastSeen, short, updateIntervalInSeconds, children, flipped } = props

  return (
    <Status status="good" pulseTrigger={lastSeen} flipped={flipped} ref={ref}>
      <div className={classnames(className, style.container)}>
        {!short && <Message className={style.message} content={sharedMessages.lastSeen} />}
        <DateTime.Relative
          value={lastSeen}
          computeDelta={computeDeltaInSeconds}
          firstToLower={!short}
          updateIntervalInSeconds={updateIntervalInSeconds}
          relativeTimeStyle={short ? 'short' : undefined}
        />
      </div>
      {children}
    </Status>
  )
})

LastSeen.propTypes = {
  children: PropTypes.node,
  className: PropTypes.string,
  flipped: PropTypes.bool,
  lastSeen: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number, // Support timestamps.
    PropTypes.instanceOf(Date),
  ]).isRequired,
  short: PropTypes.bool,
  updateIntervalInSeconds: PropTypes.number,
}

LastSeen.defaultProps = {
  children: undefined,
  className: undefined,
  flipped: false,
  updateIntervalInSeconds: undefined,
  short: false,
}

export default LastSeen
