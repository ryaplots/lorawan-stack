// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

import Link from '@ttn-lw/components/link'

import Message from '@ttn-lw/lib/components/message'

import sharedMessages from '@ttn-lw/lib/shared-messages'
import PropTypes from '@ttn-lw/lib/prop-types'

import style from './doc.styl'

import Tooltip from '.'

const DocTooltip = ({ content, docPath, docTitle, interactive, ...rest }) => (
  <Tooltip
    interactive
    content={
      <>
        {content}
        <div className={style.docLink}>
          <Link.DocLink primary path={docPath}>
            <Message content={docTitle} />
          </Link.DocLink>
        </div>
      </>
    }
    {...rest}
    appendTo={document.body}
  />
)

DocTooltip.propTypes = {
  ...Tooltip.propTypes,
  docPath: PropTypes.string.isRequired,
  docTitle: PropTypes.message,
}

DocTooltip.defaultProps = {
  ...Tooltip.defaultProps,
  docTitle: sharedMessages.moreInformation,
}

export default DocTooltip
