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
import { defineMessages } from 'react-intl'

import deviceIcon from '@assets/misc/end-device.svg'

import Status from '@ttn-lw/components/status'
import Tooltip from '@ttn-lw/components/tooltip'
import DocTooltip from '@ttn-lw/components/tooltip/doc'
import Icon from '@ttn-lw/components/icon'

import Message from '@ttn-lw/lib/components/message'

import EntityTitleSection from '@console/components/entity-title-section'
import LastSeen from '@console/components/last-seen'

import PropTypes from '@ttn-lw/lib/prop-types'
import sharedMessages from '@ttn-lw/lib/shared-messages'

import style from './device-title-section.styl'

const m = defineMessages({
  uplinkDownlinkTooltip:
    'The number of sent uplinks and received downlinks of this end device since the last frame counter reset.',
  lastSeenAvailableTooltip:
    'The elapsed time since the network registered the last activity of this end device. This is determined from sent uplinks, confirmed downlinks or (re)join requests.',
  noActivityTooltip:
    'The network has not registered any activity from this end device yet. This could mean that your end device has not sent any messages yet or only messages that cannot be handled by the network, e.g. due to a mismatch of EUIs or frequencies.',
})

const { Content } = EntityTitleSection

const DeviceTitleSection = props => {
  const {
    devId,
    fetching,
    device,
    uplinkFrameCount,
    downlinkFrameCount,
    lastSeen,
    children,
  } = props

  const showLastSeen = Boolean(lastSeen)
  const showUplinkCount = typeof uplinkFrameCount === 'number'
  const showDownlinkCount = typeof downlinkFrameCount === 'number'
  const notAvailableElem = <Message content={sharedMessages.notAvailable} />
  const bottomBarLeft = (
    <>
      <Tooltip content={<Message content={m.uplinkDownlinkTooltip} />}>
        <div className={style.messages}>
          <Content.MessagesCount
            icon="uplink"
            value={showUplinkCount ? uplinkFrameCount : notAvailableElem}
            iconClassName={showUplinkCount ? style.messageIcon : style.notAvailable}
          />
          <Content.MessagesCount
            icon="downlink"
            value={showDownlinkCount ? downlinkFrameCount : notAvailableElem}
            iconClassName={showUplinkCount ? style.messageIcon : style.notAvailable}
          />
        </div>
      </Tooltip>
      {showLastSeen ? (
        <DocTooltip
          docPath="/reference/last-activity"
          content={<Message content={m.lastSeenAvailableTooltip} />}
        >
          <LastSeen lastSeen={lastSeen} flipped>
            <Icon icon="help_outline" textPaddedLeft small nudgeUp className="tc-subtle-gray" />
          </LastSeen>
        </DocTooltip>
      ) : (
        <DocTooltip
          docPath="/devices/troubleshooting/#my-device-wont-join-what-do-i-do"
          docTitle={sharedMessages.troubleshooting}
          content={<Message content={m.noActivityTooltip} />}
        >
          <Status status="mediocre" label={sharedMessages.noActivityYet} flipped>
            <Icon icon="help_outline" textPaddedLeft small nudgeUp className="tc-subtle-gray" />
          </Status>
        </DocTooltip>
      )}
    </>
  )

  return (
    <EntityTitleSection
      id={devId}
      name={device.name}
      icon={deviceIcon}
      iconAlt={sharedMessages.device}
    >
      <Content
        className={style.content}
        creationDate={device.created_at}
        fetching={fetching}
        bottomBarLeft={bottomBarLeft}
      />
      {children}
    </EntityTitleSection>
  )
}

DeviceTitleSection.propTypes = {
  children: PropTypes.oneOfType([PropTypes.arrayOf(PropTypes.node), PropTypes.node]),
  devId: PropTypes.string.isRequired,
  device: PropTypes.device.isRequired,
  downlinkFrameCount: PropTypes.number,
  fetching: PropTypes.bool,
  lastSeen: PropTypes.string,
  uplinkFrameCount: PropTypes.number,
}

DeviceTitleSection.defaultProps = {
  uplinkFrameCount: undefined,
  lastSeen: undefined,
  children: null,
  fetching: false,
  downlinkFrameCount: undefined,
}

export default DeviceTitleSection
